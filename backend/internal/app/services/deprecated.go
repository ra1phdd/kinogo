package services

/*func (s *Service) SearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error("Ошибка парсинга из формы", zap.Error(err))
		return
	}

	textSearch := strings.ToLower(r.Form.Get("search"))
	logger.Debug("Получение текста поиска из формы", zap.String("textSearch", textSearch))

	movies, err := s.GetSearchMovies(textSearch)
	if err != nil {
		logger.Error("Ошибка получения контента по поиску", zap.Error(err), zap.Any("movies", movies))
		return
	}
	logger.Debug("Получен контент из БД/кэша", zap.Any("movies", movies))

	var sb strings.Builder
	for _, movie := range movies {
		sb.WriteString(fmt.Sprintf("<a href='/id/%d'>%s (%d)</a>", movie.Id, movie.Title, movie.ReleaseDate))
	}

	_, _ = fmt.Fprintln(w, sb.String())
}

func (s *Service) ValidateInput(input string) error {
	// Проверка на пустую строку
	if strings.TrimSpace(input) == "" {
		return errors.New("входные данные пустые")
	}

	// Проверка на максимальную длину
	const maxLength = 200
	if len(input) > maxLength {
		return errors.New("входные данные слишком длинные (больше 200 символов)")
	}

	// Проверка на специальные символы
	specialChars := []string{"<", ">", "&", "%"}
	for _, char := range specialChars {
		if strings.Contains(input, char) {
			return errors.New("входные данные содержат специальные символы (<, >, &, %)")
		}
	}

	// Все проверки пройдены, возвращаем nil
	return nil
}

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

func (s *Service) HandleView(r *http.Request, contentID int64) {
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
			logger.Error("Ошибка получения значения views из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения views на 1 и сохранение в Redis
		viewsKey := fmt.Sprintf("views:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, viewsKey, views+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipViewKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaViewKey, "1", 0)
		logger.Debug("Просмотр зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("views", views+1))
	}
}

func (s *Service) HandleLike(r *http.Request, contentID int64) {
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
			logger.Error("Ошибка получения значения likes из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения likes на 1 и сохранение в Redis
		likesKey := fmt.Sprintf("likes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, likesKey, likes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipLikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaLikeKey, "1", 0)
		logger.Debug("Лайк зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("likes", likes+1))
	}
}

func (s *Service) HandleDislike(r *http.Request, contentID int64) {
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
			logger.Error("Ошибка получения значения dislikes из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения dislikes на 1 и сохранение в Redis
		dislikesKey := fmt.Sprintf("dislikes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, dislikesKey, dislikes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipDislikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaDislikeKey, "1", 0)
		logger.Debug("Дизлайк зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("dislikes", dislikes+1))
	}
}

func (s *Service) SaveStatsToDB() {
	// Получение данных из Redis
	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о просмотрах из Redis", zap.Error(err))
		return
	}

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о лайках из Redis", zap.Error(err))
		return
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о дизлайках из Redis", zap.Error(err))
		return
	}

	logger.Debug("Данные из Redis получены", zap.Strings("views", views), zap.Strings("likes", likes), zap.Strings("dislikes", dislikes))

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		// Обновление значения views в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET views = $1 WHERE id = $2`, views, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения views в базе данных", zap.Error(err))
			return
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()
		// Обновление значения likes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET likes = $1 WHERE id = $2`, likes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения likes в базе данных", zap.Error(err))
			return
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()
		// Обновление значения dislikes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET dislikes = $1 WHERE id = $2`, dislikes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения dislikes в базе данных", zap.Error(err))
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

func (s *Service) GetStatsToDB(id int) (int64, int64, int64, error) {
	var allViews int64
	var allLikes int64
	var allDislikes int64

	logger.Debug("Получение данных из базы данных", zap.Int("id", id))

	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о просмотрах из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о лайках из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о дизлайках из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		if contentID == int64(id) {
			allViews = views
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()
		if contentID == int64(id) {
			allLikes = likes
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()
		if contentID == int64(id) {
			allDislikes = dislikes
		}
	}

	logger.Debug("Данные из Redis по фильму получены", zap.Int64("views", allViews), zap.Int64("likes", allLikes), zap.Int64("dislikes", allDislikes))

	return allViews, allLikes, allDislikes, nil
}
*/
