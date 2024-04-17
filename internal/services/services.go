package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	"kinogo/internal/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

// Функция проверки наличия стрима
func IsStreaming() (bool, error) {
	/*movies, err := rdb.Get(cache.Ctx, "isStreaming").Bool()
	if err == nil {
		return movies, err
	}

	// Данных нет в Redis, получаем их из базы данных
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login=zubarefff", nil)
	if err != nil {
		logger.Warn("Ошибка GET-запроса к API Twitch", zap.Error(err))
		return false, err
	}

	req.Header.Add("Authorization", "Bearer 6q3g8ycgdu9faa9hw1w92l1poyn0ku")
	req.Header.Add("Client-ID", "gp762nuuoqcoxypju8c569th9wz7q5")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Warn("Ошибка отправки HTTP-запроса к API Twitch", zap.Error(err))
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Warn("Ошибка чтения ответа от API Twitch", zap.Error(err))
		return false, err
	}

	var isStreaming bool
	bodyStr := string(body)
	if bodyStr != `{"data":[],"pagination":{}}` {
		isStreaming = true
	} else {
		isStreaming = false
	}

	err = rdb.Set(cache.Ctx, "isStreaming", isStreaming, 5*time.Minute).Err()
	if err != nil {
		return false, err
	}

	if isStreaming {
		return true, nil
	} else {
		return false, nil
	}*/
	return false, nil
}

// Функция вывода результатов поиска по API через AJAX в админ-панели
func ResultMovieHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	movieName := r.Form.Get("movieName")
	if err := ValidateInput(movieName); err != nil {
		logger.Warn("Ошибка валидации входных данных", err)
		return
	}

	url := "https://api.kinopoisk.dev/v1.4/movie/search?query=" + url.QueryEscape(movieName) + "&limit=4"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Warn("Ошибка GET-запроса к API КиноПоиска", err)
		return
	}

	req.Header.Add("X-API-KEY", "PNRS21P-Q0746F9-J85KRM9-S5YR004")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Warn("Ошибка отправки HTTP-запроса к API КиноПоиска", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Warn("Ошибка чтения ответа от API КиноПоиска", err)
		return
	}

	var sb strings.Builder
	err = json.Unmarshal(body, &movies)
	if err != nil {
		logger.Warn("Ошибка декодирования JSON", err)
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

	file, _, err := r.FormFile("send-video")
	if err != nil {
		logger.Error("Ошибка при парсинге видео из формы", err)
		return
	}
	defer file.Close()

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
				logger.Error("Ошибка добавления фильма в таблицу movies", err)
				return
			}

			for _, genre := range doc.Genres {
				var idGenre int
				var nameGenre string

				rows, err := db.Conn.Query("SELECT * FROM genres")
				if err != nil {
					logger.Warn("Ошибка выборки жанров из таблицы genres", err)
				}
				defer rows.Close()

				for rows.Next() {
					err := rows.Scan(&idGenre, &nameGenre)
					if err != nil {
						logger.Warn("Ошибка чтения строки из результата SQL-запроса", err)
					}
					if nameGenre == genre.Name {
						_, err := db.Conn.Exec(`INSERT INTO MoviesGenres (idmovie, idgenre) VALUES ($1, $2)`, idMovie, idGenre)
						if err != nil {
							logger.Warn("Ошибка добавления связи фильма-жанров в таблицу MoviesGenres", err)
						}
					}
				}

				if err = rows.Err(); err != nil {
					logger.Warn("Ошибка rows.Err() в функции AddMovieHandler", err)
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
				logger.Error("Ошибка создания файла фильма", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				logger.Error("Ошибка копирования файла фильма из формы на диск", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go ProcessingFile(idMovie, &wg)
			wg.Wait()
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
		fmt.Println("Ошибка при выполнении команды ffmpeg:", err)
	}

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
		fmt.Println("Ошибка при выполнении команды ffmpeg:", err)
		return
	}
	totalDuration := gjson.Get(a, "format.duration").Float()

	err = ffmpeg.Input(inFileName).
		Output(outFileName, ffmpeg.KwArgs{
			"c:v": "av1_nvenc",
			"b:v": args[1],
			"c:a": "aac",
			"b:a": "128k",
			"map": "0",
			"s":   args[2],
			"f":   "dash",
		}).GlobalArgs("-progress", "unix://"+TempSock(totalDuration, progress)).
		OverWriteOutput().
		Run()

	if err != nil {
		fmt.Println("Ошибка при выполнении команды ffmpeg:", err)
	}
}

func TempSock(totalDuration float64, progress *Progress) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock_%d", r.Int(), time.Now().UnixNano()))
	os.Remove(sockFileName)
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
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

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	textSearch := strings.ToLower(r.Form.Get("search"))

	movies, err := GetSearchMovies(textSearch)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	fmt.Println(movies)

	var sb strings.Builder
	for _, movie := range movies {
		sb.WriteString(fmt.Sprintf("<a href='/id/%d'>%s (%d)</a>", movie.Id, movie.Title, movie.ReleaseDate))
	}

	fmt.Fprintln(w, sb.String())
}

func ValidateInput(input string) error {
	// Проверка на пустую строку
	if strings.TrimSpace(input) == "" {
		return errors.New("input cannot be empty")
	}

	// Проверка на максимальную длину
	const maxLength = 200
	if len(input) > maxLength {
		return errors.New("input is too long")
	}

	// Проверка на специальные символы
	specialChars := []string{"<", ">", "&", "%"}
	for _, char := range specialChars {
		if strings.Contains(input, char) {
			return errors.New("input contains invalid characters")
		}
	}

	// Все проверки пройдены, возвращаем nil
	return nil
}

func GetAllMovies() ([]models.MovieData, error) {
	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllMovies").Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
        FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
        JOIN genres ON moviesgenres.idgenre = genres.id GROUP BY movies.id ORDER BY movies.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке данных из таблицы arrays: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedulessss: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		fmt.Println(MoviesData)

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllMovies", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func GetBestMovie() (models.MovieData, error) {
	var moviesSlice models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryBestMovie").Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return moviesSlice, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres FROM movies
		JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id
		GROUP BY movies.id ORDER BY movies.views DESC LIMIT 1`)
	if err != nil {
		return moviesSlice, fmt.Errorf("ошибка при выборке данных из таблицы arrays: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return moviesSlice, fmt.Errorf("ошибка при обработке данных из таблицы schedulessss: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		fmt.Println(MoviesData)

		moviesSlice = MoviesData
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return moviesSlice, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryBestMovie", moviesJSON, 5*time.Minute).Err()
	if err != nil {
		return moviesSlice, err
	}

	return moviesSlice, nil
}

func GetFilterMovies(idGenres []int, yearMin int, yearMax int) ([]models.MovieData, error) {
	var idGenresStr []string
	for _, id := range idGenres {
		idGenresStr = append(idGenresStr, strconv.Itoa(id))
	}

	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryFilterMovies_"+strings.Join(idGenresStr, "_")+"_"+fmt.Sprint(yearMin)+"_"+fmt.Sprint(yearMax)).Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	var rows *sqlx.Rows
	if len(idGenres) > 0 {
		rows, err = db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
			FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
			JOIN genres ON moviesgenres.idgenre = genres.id
			WHERE genres.id = ANY($1) AND (releasedate >= $2 AND releasedate <= $3)
			GROUP BY movies.id ORDER BY movies.id DESC`, idGenres, yearMin, yearMax)
		if err != nil {
			fmt.Println("ошибка при выборке данных из таблицы arrays в функци и getNextSched: ", err)
		}
		defer rows.Close()
	} else {
		rows, err = db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
			FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
			JOIN genres ON moviesgenres.idgenre = genres.id WHERE
			(releasedate >= $1 AND releasedate <= $2) GROUP BY movies.id ORDER BY movies.id DESC`, yearMin, yearMax)
		if err != nil {
			fmt.Println("ошибка при выборке данных из таблицы arrays в функции getNextSched: ", err)
		}
		defer rows.Close()
	}

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedules: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryFilterMovies_"+strings.Join(idGenresStr, "_")+"_"+fmt.Sprint(yearMin)+"_"+fmt.Sprint(yearMax), moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func GetMovie(id int) (models.MovieData, error) {
	var moviesSlice models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryMovie_"+fmt.Sprint(id)).Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return moviesSlice, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*,array_agg(genres.name) AS genres
		FROM movies	JOIN moviesgenres ON movies.id = moviesgenres.idmovie
		JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.id = $1 GROUP BY movies.id`, id)
	if err != nil {
		fmt.Println("ошибка при выборке данных из таблицы arrays в функции getNextSched: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&moviesSlice.Id,
			&moviesSlice.Title,
			&moviesSlice.Description,
			&moviesSlice.Country,
			&moviesSlice.ReleaseDate,
			&moviesSlice.TimeMovie,
			&moviesSlice.ScoreKP,
			&moviesSlice.ScoreIMDB,
			&moviesSlice.Poster,
			&moviesSlice.TypeMovie,
			&moviesSlice.Views,
			&moviesSlice.Likes,
			&moviesSlice.Dislikes,
			&moviesSlice.Genres,
		)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}
	}

	moviesSlice.Genres = strings.Replace(strings.Trim(moviesSlice.Genres, "{}"), ",", ", ", -1)

	moviesSlice.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", moviesSlice.ScoreKP), 64)
	if err != nil {
		fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
	}

	moviesSlice.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", moviesSlice.ScoreIMDB), 64)
	if err != nil {
		fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return moviesSlice, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryMovie_"+fmt.Sprint(id), moviesJSON, 24*time.Hour).Err()
	if err != nil {
		return moviesSlice, err
	}

	return moviesSlice, nil
}

func GetSearchMovies(textSearch string) ([]models.MovieData, error) {
	fmt.Println(textSearch)
	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QuerySearchMovies_"+textSearch).Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres FROM movies
		JOIN moviesgenres ON movies.id = moviesgenres.idmovie JOIN genres ON moviesgenres.idgenre = genres.id
		WHERE word_similarity(movies.title, $1) > 0.1 GROUP BY movies.id`, textSearch)
	if err != nil {
		fmt.Println("ошибка при выборке данных из таблицы arrays в функции getNextSched: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedules: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QuerySearchMovies_"+textSearch, moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func GetAllFilms() ([]models.MovieData, error) {
	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllFilms").Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
        FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
        JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'movie' GROUP BY movies.id ORDER BY movies.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке данных из таблицы arrays: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedulessss: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		fmt.Println(MoviesData)

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllFilms", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func GetAllCartoons() ([]models.MovieData, error) {
	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllCartoons").Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
        FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
        JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'cartoon' GROUP BY movies.id ORDER BY movies.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке данных из таблицы arrays: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedulessss: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		fmt.Println(MoviesData)

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllCartoons", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func GetAllTelecasts() ([]models.MovieData, error) {
	var moviesSlice []models.MovieData
	movies, err := cache.Rdb.Get(cache.Ctx, "QueryAllTelecasts").Result()
	if err == nil && movies != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(movies), &moviesSlice)
		if err != nil {
			return nil, err
		}
		return moviesSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT movies.*, array_agg(genres.name) AS genres
        FROM movies JOIN moviesgenres ON movies.id = moviesgenres.idmovie
        JOIN genres ON moviesgenres.idgenre = genres.id WHERE movies.typemovie = 'telecast' GROUP BY movies.id ORDER BY movies.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке данных из таблицы arrays: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var MoviesData models.MovieData
		err := rows.Scan(
			&MoviesData.Id,
			&MoviesData.Title,
			&MoviesData.Description,
			&MoviesData.Country,
			&MoviesData.ReleaseDate,
			&MoviesData.TimeMovie,
			&MoviesData.ScoreKP,
			&MoviesData.ScoreIMDB,
			&MoviesData.Poster,
			&MoviesData.TypeMovie,
			&MoviesData.Views,
			&MoviesData.Likes,
			&MoviesData.Dislikes,
			&MoviesData.Genres,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке данных из таблицы schedulessss: %w", err)
		}

		MoviesData.Genres = strings.Replace(strings.Trim(MoviesData.Genres, "{}"), ",", ", ", -1)

		MoviesData.ScoreKP, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreKP), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		MoviesData.ScoreIMDB, err = strconv.ParseFloat(fmt.Sprintf("%.1f", MoviesData.ScoreIMDB), 64)
		if err != nil {
			fmt.Println("ошибка при обработке данных из таблицы schedules в функции getNextSchedule: %w", err)
		}

		fmt.Println(MoviesData)

		moviesSlice = append(moviesSlice, MoviesData)
	}

	// Сохраняем данные в Redis
	moviesJSON, err := json.Marshal(moviesSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryAllTelecasts", moviesJSON, 1*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return moviesSlice, nil
}

func HandleView(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как просмотревший данный контент
	ipViewKey := fmt.Sprintf("ip_view:%d:%s", contentID, ipAddress)
	uaViewKey := fmt.Sprintf("ua_view:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик просмотров и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipViewKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaViewKey).Val() == 0 {
		// Получение текущего значения views из базы данных
		var views int64
		err := db.Conn.QueryRow(`SELECT views FROM movies WHERE id = $1`, contentID).Scan(&views)
		if err != nil {
			logger.Error("Ошибка получения значения views из базы данных", err)
			return
		}

		// Увеличение значения views на 1 и сохранение в Redis
		viewsKey := fmt.Sprintf("views:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, viewsKey, views+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipViewKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaViewKey, "1", 0)
		fmt.Println("Просмотр зарегистрирован. Общее количество просмотров: ", views+1)
	} else {
		fmt.Println("Вы уже просматривали этот контент")
	}
}

func HandleLike(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как оценивший данный контент
	ipLikeKey := fmt.Sprintf("ip_like:%d:%s", contentID, ipAddress)
	uaLikeKey := fmt.Sprintf("ua_like:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик лайков и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipLikeKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaLikeKey).Val() == 0 {
		// Получение текущего значения likes из базы данных
		var likes int64
		err := db.Conn.QueryRow(`SELECT likes FROM movies WHERE id = $1`, contentID).Scan(&likes)
		if err != nil {
			logger.Error("Ошибка получения значения likes из базы данных", err)
			return
		}

		// Увеличение значения likes на 1 и сохранение в Redis
		likesKey := fmt.Sprintf("likes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, likesKey, likes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipLikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaLikeKey, "1", 0)
		fmt.Println("Лайк зарегистрирован. Общее количество лайков: ", likes+1)
	} else {
		fmt.Println("Вы уже оценивали этот контент")
	}
}

func HandleDislike(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как оценивший данный контент
	ipDislikeKey := fmt.Sprintf("ip_dislike:%d:%s", contentID, ipAddress)
	uaDislikeKey := fmt.Sprintf("ua_dislike:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик дизлайков и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipDislikeKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaDislikeKey).Val() == 0 {
		// Получение текущего значения dislikes из базы данных
		var dislikes int64
		err := db.Conn.QueryRow(`SELECT dislikes FROM movies WHERE id = $1`, contentID).Scan(&dislikes)
		if err != nil {
			logger.Error("Ошибка получения значения dislikes из базы данных", err)
			return
		}

		// Увеличение значения dislikes на 1 и сохранение в Redis
		dislikesKey := fmt.Sprintf("dislikes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, dislikesKey, dislikes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipDislikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaDislikeKey, "1", 0)
		fmt.Println("Дизлайк зарегистрирован. Общее количество дизлайков: ", dislikes+1)
	} else {
		fmt.Println("Вы уже оценивали этот контент")
	}
}

func SaveStatsToDB() {
	// Получение данных из Redis
	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		// Обработка ошибки
		return
	}

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		// Обработка ошибки
		return
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		// Обработка ошибки
		return
	}
	fmt.Println(views, likes, dislikes)

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		// Обновление значения views в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET views = $1 WHERE id = $2`, views, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения views в базе данных", err)
			return
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()
		// Обновление значения likes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET likes = $1 WHERE id = $2`, likes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения likes в базе данных", err)
			return
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()
		// Обновление значения dislikes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET dislikes = $1 WHERE id = $2`, dislikes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения dislikes в базе данных", err)
			return
		}
	}

	// Очистка данных из Redis
	for _, viewKey := range views {
		cache.Rdb.Del(cache.Ctx, viewKey)
	}

	for _, likeKey := range likes {
		cache.Rdb.Del(cache.Ctx, likeKey)
	}

	for _, dislikeKey := range dislikes {
		cache.Rdb.Del(cache.Ctx, dislikeKey)
	}
}

func GetStatsToDB(id int) (int64, int64, int64, error) {
	var allViews int64
	var allLikes int64
	var allDislikes int64

	fmt.Println(id)

	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		// Обработка ошибки
		fmt.Println("ошибка при получении данных из Redis: ", err)
	}

	fmt.Println(views)

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		// Обработка ошибки
		fmt.Println("ошибка при получении данных из Redis: ", err)
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		// Обработка ошибки
		fmt.Println("ошибка при получении данных из Redis: ", err)
	}

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		fmt.Println(views)
		if contentID == int64(id) {
			allViews = views
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()

		fmt.Println(likes)
		if contentID == int64(id) {
			allLikes = likes
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()

		fmt.Println(dislikes)
		if contentID == int64(id) {
			allDislikes = dislikes
		}
	}

	return allViews, allLikes, allDislikes, nil
}
