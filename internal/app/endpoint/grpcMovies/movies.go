package grpcMovies

import (
	"context"
	"errors"
	"kinogo/internal/app/models"
	pb "kinogo/pkg/movies_v1"
)

type Movies interface {
	GetMoviesService() ([]models.Movies, error)
	GetMovieByIdService(int32) (models.Movie, error)
	GetMoviesByFilterService(map[string]interface{}) ([]models.Movies, error)
	AddMoviesService(map[string]interface{}) (int32, error)
	DelMoviesService(int32) error
}

type Endpoint struct {
	Movies Movies
	pb.UnimplementedMoviesV1Server
}

func (e *Endpoint) GetMovies(_ context.Context, _ *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	movies, err := e.Movies.GetMoviesService()
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

func (e *Endpoint) GetMoviesById(_ context.Context, req *pb.GetMoviesByIdRequest) (*pb.GetMoviesByIdResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id новости не указан")
	}

	movie, err := e.Movies.GetMovieByIdService(req.Id)
	if err != nil {
		return &pb.GetMoviesByIdResponse{}, err
	}

	pbMovie := &pb.GetMoviesByIdResponse{
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

	return pbMovie, nil
}

func (e *Endpoint) GetMoviesByFilter(_ context.Context, req *pb.GetMoviesByFilterRequest) (*pb.GetMoviesResponse, error) {
	filtersMap := map[string]interface{}{
		"typeMovie": req.Filters.TypeMovie,
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
		"country":     req.Country,
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

func (e *Endpoint) DelMovies(_ context.Context, req *pb.DelMoviesRequest) (*pb.DelMoviesResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id новости не указан")
	}

	err := e.Movies.DelMoviesService(req.Id)
	if err != nil {
		return &pb.DelMoviesResponse{Err: error.Error(err)}, err
	}

	return &pb.DelMoviesResponse{Err: ""}, nil
}
