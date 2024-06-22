package models

import "time"

type (
	Movies struct {
		Id          int32   `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ReleaseDate int32   `json:"release_date"`
		ScoreKP     float64 `json:"score_kp"`
		ScoreIMDB   float64 `json:"score_imdb"`
		Poster      string  `json:"poster"`
		TypeMovie   int32   `json:"type_movie"`
		Genres      string  `json:"genres"`
	}

	Movie struct {
		Id          int32
		Title       string
		Description string
		Country     string
		ReleaseDate int32
		TimeMovie   int32
		ScoreKP     float64
		ScoreIMDB   float64
		Poster      string
		TypeMovie   int32
		Views       int32
		Likes       int32
		Dislikes    int32
		Genres      string
	}

	MovieByAPI struct {
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

	User struct {
		ID        int64
		FirstName string
		LastName  string
		Username  string
		PhotoURL  string
		AuthDate  int64
		Hash      string
	}

	Comments struct {
		ID        int32
		UserID    int32
		ParentID  int32
		Text      string
		CreatedAt time.Time
		UpdatedAt time.Time
		Children  []Comments
	}
)
