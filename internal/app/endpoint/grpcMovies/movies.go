package grpcMovies

import (
	"context"
	"errors"
	"kinogo/internal/app/models"
	pb "kinogo/pkg/movies_v1"
)

type Movies interface {
	GetMoviesService() ([]models.Movies, error)
	GetMoviesByIdService(int32) (models.Movies, error)
	GetMoviesByCategoryService(string) ([]models.Movies, error)
	AddMoviesService(string, string, string, string) (int32, error)
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

	return &pb.GetMoviesByIdResponse{}, nil
}

func (e *Endpoint) GetMoviesByCategory(_ context.Context, req *pb.GetMoviesByCategoryRequest) (*pb.GetMoviesResponse, error) {
	if req.Categories == "" {
		return nil, errors.New("id категории не указан")
	}

	return &pb.GetMoviesResponse{}, nil
}

func (e *Endpoint) AddMovies(_ context.Context, req *pb.AddMoviesRequest) (*pb.AddMoviesResponse, error) {
	if req.Title == "" {
		return nil, errors.New("заголовок новости не указан")
	}
	if req.Text == "" {
		return nil, errors.New("текст новости не указан")
	}
	if req.Datetime == "" {
		return nil, errors.New("дата публикации новости не указан")
	}
	if req.Categories == "" {
		return nil, errors.New("id категорий новости не указан")
	}

	id, err := e.Movies.AddMoviesService(req.Title, req.Text, req.Datetime, req.Categories)
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
