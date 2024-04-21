package models

type (
	GeneralData struct {
		Stream         bool
		Auth           bool
		TextSearch     string
		IndexHandler   bool
		MovieHandler   bool
		SearchHandler  bool
		FilterHandler  bool
		SearchAside    bool
		FilterAside    bool
		BestMovieAside bool
	}

	MovieData struct {
		Id          int64
		Title       string
		Description string
		Country     string
		ReleaseDate int
		TimeMovie   int
		ScoreKP     float64
		ScoreIMDB   float64
		Poster      string
		TypeMovie   string
		Views       int64
		Likes       int64
		Dislikes    int64
		Genres      string
	}

	FilterData struct {
		BoolFilter bool
		Genre      []string
		YearMin    int
		YearMax    int
	}

	AllData struct {
		GeneralData   GeneralData
		MovieData     []MovieData
		FilterData    FilterData
		BestMovieData MovieData
		UserData      User
		CommentsData  []Comment
	}

	Movies struct {
		Docs []struct {
			ID                int    `json:"id"`
			Name              string `json:"name"`
			AlternativeName   string `json:"alternativeName"`
			EnName            string `json:"enName"`
			Type              string `json:"type"`
			Year              int    `json:"year"`
			Description       string `json:"description"`
			ShortDescription  string `json:"shortDescription"`
			MovieLength       int    `json:"movieLength"`
			IsSeries          bool   `json:"isSeries"`
			TicketsOnSale     bool   `json:"ticketsOnSale"`
			TotalSeriesLength any    `json:"totalSeriesLength"`
			SeriesLength      any    `json:"seriesLength"`
			RatingMpaa        string `json:"ratingMpaa"`
			AgeRating         int    `json:"ageRating"`
			Top10             any    `json:"top10"`
			Top250            any    `json:"top250"`
			TypeNumber        int    `json:"typeNumber"`
			Status            any    `json:"status"`
			Names             []struct {
				Name     string `json:"name"`
				Language string `json:"language,omitempty"`
				Type     any    `json:"type,omitempty"`
			} `json:"names"`
			ExternalID struct {
				Imdb string `json:"imdb"`
				Tmdb int    `json:"tmdb"`
				KpHD string `json:"kpHD"`
			} `json:"externalId"`
			Logo struct {
				URL string `json:"url"`
			} `json:"logo"`
			Poster struct {
				URL        string `json:"url"`
				PreviewURL string `json:"previewUrl"`
			} `json:"poster"`
			Backdrop struct {
				URL        string `json:"url"`
				PreviewURL string `json:"previewUrl"`
			} `json:"backdrop"`
			Rating struct {
				Kp                 float64 `json:"kp"`
				Imdb               float64 `json:"imdb"`
				FilmCritics        float64 `json:"filmCritics"`
				RussianFilmCritics float64 `json:"russianFilmCritics"`
				Await              any     `json:"await"`
			} `json:"rating"`
			Votes struct {
				Kp                 int `json:"kp"`
				Imdb               int `json:"imdb"`
				FilmCritics        int `json:"filmCritics"`
				RussianFilmCritics int `json:"russianFilmCritics"`
				Await              int `json:"await"`
			} `json:"votes"`
			Genres []struct {
				Name string `json:"name"`
			} `json:"genres"`
			Countries []struct {
				Name string `json:"name"`
			} `json:"countries"`
			ReleaseYears []any `json:"releaseYears"`
		} `json:"docs"`
		Total int `json:"total"`
		Limit int `json:"limit"`
		Page  int `json:"page"`
		Pages int `json:"pages"`
	}

	// Структура данных с информацией о пользователе
	User struct {
		ID        int64
		FirstName string
		LastName  string
		Username  string
		PhotoURL  string
		AuthDate  int64
		Hash      string
	}

	Comment struct {
		ID       int
		ParentID int
		Text     string
		MovieID  int
		User     User
		Children []*Comment
	}
)
