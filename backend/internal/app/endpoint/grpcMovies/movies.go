package grpcMovies

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"kinogo/internal/app/models"
	"kinogo/pkg/logger"
	pb "kinogo/pkg/movies_v1"
	"strings"
)

type Movies interface {
	GetMoviesService(int32, int32) ([]models.Movies, error)
	GetMovieByIdService(int32) (models.Movie, error)
	GetMoviesByFilterService(map[string]interface{}) ([]models.Movies, error)
	AddMoviesService(map[string]interface{}) (int32, error)
	DeleteMoviesService(int32) error
}

type Endpoint struct {
	Movies Movies
	pb.UnimplementedMoviesV1Server
}

func (e *Endpoint) GetMovies(_ context.Context, req *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	movies, err := e.Movies.GetMoviesService(req.Limit, req.Page)
	if err != nil {
		logger.Error("Ошибка в работе функции GetMoviesService", zap.Error(err))
		return &pb.GetMoviesResponse{}, err
	}

	var pbMovies []*pb.GetMovieItem
	for _, movie := range movies {
		pbMovie := &pb.GetMovieItem{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			ReleaseDate: movie.ReleaseDate,
			ScoreKP:     movie.ScoreKP,
			ScoreIMDB:   movie.ScoreIMDB,
			Poster:      movie.Poster,
			TypeMovie:   movie.TypeMovie,
			Genres:      movie.Genres,
		}
		pbMovies = append(pbMovies, pbMovie)
	}

	return &pb.GetMoviesResponse{Movies: pbMovies}, nil
}

func (e *Endpoint) GetMovieById(_ context.Context, req *pb.GetMoviesByIdRequest) (*pb.GetMoviesByIdResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id новости не указан")
	}

	movie, err := e.Movies.GetMovieByIdService(req.Id)
	if err != nil {
		if !strings.Contains(err.Error(), "NotFound") {
			logger.Error("Ошибка в работе функции GetMovieByIdService", zap.Error(err))
		}
		return &pb.GetMoviesByIdResponse{}, err
	}

	pbMovie := &pb.GetMoviesByIdResponse{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		Country:     movie.Country,
		ReleaseDate: movie.ReleaseDate,
		TimeMovie:   movie.TimeMovie,
		ScoreKP:     movie.ScoreKP,
		ScoreIMDB:   movie.ScoreIMDB,
		Poster:      movie.Poster,
		TypeMovie:   movie.TypeMovie,
		Genres:      movie.Genres,
	}

	return pbMovie, nil
}

func (e *Endpoint) GetMoviesByFilter(_ context.Context, req *pb.GetMoviesByFilterRequest) (*pb.GetMoviesResponse, error) {
	filtersMap := map[string]interface{}{
		"typeMovie": req.Filters.TypeMovie,
		"search":    req.Filters.Search,
		"genres":    req.Filters.Genres,
		"yearMin":   req.Filters.YearMin,
		"yearMax":   req.Filters.YearMax,
		"bestMovie": req.Filters.BestMovie,
		"limit":     req.Limit,
		"page":      req.Page,
	}

	movies, err := e.Movies.GetMoviesByFilterService(filtersMap)
	if err != nil {
		return &pb.GetMoviesResponse{}, err
	}

	var pbMovies []*pb.GetMovieItem
	for _, movie := range movies {
		pbMovie := &pb.GetMovieItem{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			ReleaseDate: movie.ReleaseDate,
			ScoreKP:     movie.ScoreKP,
			ScoreIMDB:   movie.ScoreIMDB,
			Poster:      movie.Poster,
			TypeMovie:   movie.TypeMovie,
			Genres:      movie.Genres,
		}
		pbMovies = append(pbMovies, pbMovie)
	}

	return &pb.GetMoviesResponse{Movies: pbMovies}, nil
}

func (e *Endpoint) AddMovies(_ context.Context, req *pb.AddMoviesRequest) (*pb.AddMoviesResponse, error) {
	moviesMap := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"countries":   req.Countries,
		"releaseDate": req.ReleaseDate,
		"timeMovie":   req.TimeMovie,
		"scoreKP":     req.ScoreKP,
		"scoreIMDB":   req.ScoreIMDB,
		"poster":      req.Poster,
		"typeMovie":   req.TypeMovie,
		"genres":      req.Genres,
	}

	id, err := e.Movies.AddMoviesService(moviesMap)
	if err != nil {
		return &pb.AddMoviesResponse{Err: error.Error(err)}, err
	}

	return &pb.AddMoviesResponse{Id: id}, nil
}

func (e *Endpoint) DeleteMovies(_ context.Context, req *pb.DeleteMoviesRequest) (*pb.DeleteMoviesResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id новости не указан")
	}

	err := e.Movies.DeleteMoviesService(req.Id)
	if err != nil {
		return &pb.DeleteMoviesResponse{Err: error.Error(err)}, err
	}

	return &pb.DeleteMoviesResponse{Err: ""}, nil
}
