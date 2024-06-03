package services

import (
	"encoding/json"
	"fmt"
	"io"
	"kinogo/internal/app/models"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"kinogo/cmd/websocket"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"

	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

type Progress struct {
	sync.Mutex
	value float64
}

func (p *Progress) Add(v float64) {
	p.Lock()
	defer p.Unlock()
	p.value = v
}

func (p *Progress) Value() float64 {
	p.Lock()
	defer p.Unlock()
	return p.value
}

var progresses []*Progress
var movies models.Movies

// Функция вывода результатов поиска по API через AJAX в админ-панели
func ResultMovieHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	movieName := r.Form.Get("movieName")
	if err := ValidateInput(movieName); err != nil {
		logger.Error("Ошибка валидации входных данных", zap.Error(err), zap.String("movieName", movieName))
		return
	}
	logger.Debug("Получение названия фильма из формы", zap.String("movieName", movieName))

	url := "https://api.kinopoisk.dev/v1.4/movie/search?query=" + url.QueryEscape(movieName) + "&limit=4"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Ошибка GET-запроса к API КиноПоиска", zap.Error(err))
		return
	}

	req.Header.Add("X-API-KEY", "PNRS21P-Q0746F9-J85KRM9-S5YR004")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Ошибка отправки HTTP-запроса к API КиноПоиска", zap.Error(err))
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Ошибка чтения ответа от API КиноПоиска", zap.Error(err))
		return
	}

	var sb strings.Builder
	err = json.Unmarshal(body, &movies)
	if err != nil {
		logger.Error("Ошибка декодирования JSON", zap.Error(err))
		return
	}

	countResult := 0
	for _, doc := range movies.Docs {
		if doc.Name != "" && doc.Name != "null" {
			sb.WriteString(fmt.Sprintf("<input type='radio' name='movie' id='%d' value='%d'></input><label for='%d'>%s (%d)</label>", doc.ID, doc.ID, doc.ID, doc.Name, doc.Year))
			countResult += 1
		}
	}
	if countResult == 0 {
		sb.WriteString("<p class='error_notFound'>Ничего не найдено</p>")
	}

	fmt.Fprintln(w, sb.String())
}

// Функция добавления фильма из админ-панели
func AddMovieHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	radioButtonValue := r.URL.Query().Get("id")
	logger.Debug("Получение id фильма из формы", zap.String("radioButtonValue", radioButtonValue))

	file, _, err := r.FormFile("send-video")
	if err != nil {
		logger.Error("Ошибка при парсинге видео из формы", zap.Error(err))
		return
	}
	defer file.Close()
	logger.Debug("Парсинг видео из формы завершен")

	for _, doc := range movies.Docs {
		ID := fmt.Sprint(doc.ID)

		if radioButtonValue == ID {
			title := doc.Name
			description := doc.Description
			releaseDate := doc.Year
			timeMovie := doc.MovieLength
			scoreKP := doc.Rating.Kp
			scoreIMDb := doc.Rating.Imdb
			poster := doc.Poster.URL
			typeMovie := doc.Type

			var countries []string
			for _, country := range doc.Countries {
				countries = append(countries, country.Name)
			}
			country := strings.Join(countries, ", ")

			var idMovie int
			err := db.Conn.QueryRow(`INSERT INTO movies (title, description, country, releasedate, timemovie, scorekp, scoreimdb, poster, typemovie) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, title, description, country, releaseDate, timeMovie, scoreKP, scoreIMDb, poster, typeMovie).Scan(&idMovie)
			if err != nil {
				logger.Error("Ошибка добавления фильма в таблицу movies", zap.Error(err))
				return
			}
			logger.Debug("Добавление фильма в таблицу movies", zap.Int("idMovie", idMovie), zap.String("title", title),
				zap.String("description", description), zap.String("country", country), zap.Int("releasedate", releaseDate),
				zap.Int("timemovie", timeMovie), zap.Float64("scorekp", scoreKP), zap.Float64("scoreimdb", scoreIMDb),
				zap.String("poster", poster), zap.String("typemovie", typeMovie))

			for _, genre := range doc.Genres {
				var idGenre int
				var nameGenre string

				rows, err := db.Conn.Query("SELECT * FROM genres")
				if err != nil {
					logger.Error("Ошибка выборки жанров из таблицы genres", zap.Error(err))
				}
				defer rows.Close()

				for rows.Next() {
					err := rows.Scan(&idGenre, &nameGenre)
					if err != nil {
						logger.Error("Ошибка чтения строки из результата SQL-запроса", zap.Error(err))
					}
					if nameGenre == genre.Name {
						_, err := db.Conn.Exec(`INSERT INTO MoviesGenres (idmovie, idgenre) VALUES ($1, $2)`, idMovie, idGenre)
						if err != nil {
							logger.Error("Ошибка добавления связи фильма-жанров в таблицу MoviesGenres", zap.Error(err))
							return
						}
						logger.Debug("Добавление связи фильма-жанра", zap.Int("idMovie", idMovie), zap.Int("idGenre", idGenre))
					}
				}

				if err = rows.Err(); err != nil {
					logger.Error("Ошибка rows.Err() в функции AddMovieHandler", zap.Error(err))
					return
				}
			}

			dirPath := "media/" + fmt.Sprint(idMovie)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				os.MkdirAll(dirPath, 0755)
				os.MkdirAll(dirPath+"/1080p", 0755)
				os.MkdirAll(dirPath+"/720p", 0755)
				os.MkdirAll(dirPath+"/480p", 0755)
				os.MkdirAll(dirPath+"/360p", 0755)
			}

			filePath := dirPath + "/" + fmt.Sprint(idMovie) + ".mp4"
			out, err := os.Create(filePath)
			if err != nil {
				logger.Error("Ошибка создания файла фильма", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				logger.Error("Ошибка копирования файла фильма из формы на диск", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			logger.Debug("Копирование файла фильма из формы на диск завершено", zap.Int("idMovie", idMovie))

			var wg sync.WaitGroup
			wg.Add(1)
			go ProcessingFile(idMovie, &wg)
			wg.Wait()
			logger.Debug("Обработка файла фильма завершена", zap.Int("idMovie", idMovie))
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ProcessingFile(idMovie int, wg *sync.WaitGroup) {
	err := ffmpeg.Input(fmt.Sprintf("media/%s/%s.mp4", fmt.Sprint(idMovie), fmt.Sprint(idMovie))).
		Output(fmt.Sprintf("media/%s/%s.mkv", fmt.Sprint(idMovie), fmt.Sprint(idMovie)), ffmpeg.KwArgs{"c:v": "copy", "c:a": "copy"}).
		GlobalArgs().
		OverWriteOutput().
		Run()
	if err != nil {
		logger.Error("Ошибка при выполнении команды ffmpeg для преобразования mp4 в mkv", zap.Error(err))
		return
	}
	logger.Debug("Преобразование mp4 в mkv завершено", zap.Int("idMovie", idMovie))

	os.Remove(fmt.Sprintf("media/%s/%s.mp4", fmt.Sprint(idMovie), fmt.Sprint(idMovie)))

	commands := [][]string{
		{
			fmt.Sprintf("media/%s/%s.mkv", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
			"8M",
			"hd1080",
			fmt.Sprintf("media/%s/1080p/%s_1080p.mpd", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
		},
		{
			fmt.Sprintf("media/%s/%s.mkv", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
			"6M",
			"hd720",
			fmt.Sprintf("media/%s/720p/%s_720p.mpd", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
		},
		{
			fmt.Sprintf("media/%s/%s.mkv", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
			"4M",
			"854x480",
			fmt.Sprintf("media/%s/480p/%s_480p.mpd", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
		},
		{
			fmt.Sprintf("media/%s/%s.mkv", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
			"2M",
			"640x360",
			fmt.Sprintf("media/%s/360p/%s_360p.mpd", fmt.Sprint(idMovie), fmt.Sprint(idMovie)),
		},
	}

	progresses = make([]*Progress, len(commands))
	for i := range progresses {
		progresses[i] = &Progress{}
	}

	logger.Debug("Запуск обработки видео на разных качествах", zap.Int("idMovie", idMovie))

	wg.Add(len(commands))
	for i, cmdArgs := range commands {
		go ExecuteCommand(cmdArgs, progresses[i], wg)
	}
	go PrintProgress(progresses)
}

func ExecuteCommand(args []string, progress *Progress, wg *sync.WaitGroup) {
	defer wg.Done()

	inFileName := args[0]
	outFileName := args[3]

	a, err := ffmpeg.Probe(inFileName)
	if err != nil {
		logger.Error("Ошибка при получении метаданных о файле", zap.Error(err), zap.String("inFileName", inFileName))
		return
	}
	totalDuration := gjson.Get(a, "format.duration").Float()

	sockFileName := TempSock(totalDuration, progress)
	if sockFileName == "" {
		return
	}

	err = ffmpeg.Input(inFileName).
		Output(outFileName, ffmpeg.KwArgs{
			"c:v": "av1_nvenc",
			"b:v": args[1],
			"c:a": "aac",
			"b:a": "128k",
			"map": "0",
			"s":   args[2],
			"f":   "dash",
		}).GlobalArgs("-progress", "unix://"+sockFileName).
		OverWriteOutput().
		Run()
	if err != nil {
		logger.Error("Ошибка при выполнении команды ffmpeg", zap.Error(err))
	}
}

func TempSock(totalDuration float64, progress *Progress) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock_%d", r.Int(), time.Now().UnixNano()))
	os.Remove(sockFileName)
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		logger.Error("Ошибка создания сокета", zap.Error(err), zap.String("sockFileName", sockFileName))
		return ""
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			logger.Error("Ошибка принятия входящего соединения на сокете", zap.Error(err), zap.String("sockFileName", sockFileName))
			return
		}
		buf := make([]byte, 16)
		data := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				logger.Error("Ошибка чтения данных из сокета", zap.Error(err), zap.String("sockFileName", sockFileName))
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := 0.0
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = float64(c) / totalDuration / 1000000
			}

			if strings.Contains(data, "progress=end") {
				cp = 1.0
			}

			progress.Add(cp)
		}
	}()

	return sockFileName
}

func PrintProgress(progresses []*Progress) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		totalProgress := 0.0
		for _, p := range progresses {
			totalProgress += p.Value()
		}
		if len(progresses) > 0 {
			averageProgress := totalProgress / float64(len(progresses))
			progress := fmt.Sprintf("%.3f", averageProgress)
			if websocket.Conn != nil {
				websocket.Conn.WriteJSON(map[string]string{"progress": progress})
			}
			if progress == "100" {
				return
			}
		}

	}
}

func (s *Service) ParseTemplatesAdmin(w http.ResponseWriter, allData models.AllData) {
	tmpl, err := template.ParseFiles(
		"web/admin/templates/index.html",
	)
	if err != nil {
		logger.Error("Ошибка парсинга шаблонов", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, allData)
	if err != nil {
		logger.Error("Ошибка выполнения шаблонов", zap.Error(err), zap.Any("allData", allData))
		return
	}
}
