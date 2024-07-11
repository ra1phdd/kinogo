package grpcComments

import (
	"context"
	"errors"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"kinogo/internal/app/models"
	metrics "kinogo/internal/app/services/metrics"
	pb "kinogo/pkg/comments_v1"
)

type Comments interface {
	GetCommentsByIdService(movieId int32, limit int32, page int32) ([]models.Comments, error)
	AddCommentService(data map[string]interface{}) (int32, error)
	UpdateCommentService(data map[string]interface{}) error
	DelCommentService(id int32) error
}

type Endpoint struct {
	Comments Comments
	pb.UnimplementedCommentsV1Server
}

func (e *Endpoint) GetCommentsById(_ context.Context, req *pb.GetCommentsByIdRequest) (*pb.GetCommentsByIdResponse, error) {
	if req.MovieId == 0 {
		return nil, errors.New("id фильма не указан")
	}

	comments, err := e.Comments.GetCommentsByIdService(req.MovieId, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	pbComments := convertCommentsToPb(comments)

	return &pb.GetCommentsByIdResponse{Comments: pbComments}, nil
}

func (e *Endpoint) AddComment(_ context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	commentMap := map[string]interface{}{
		"parentId":  req.ParentId,
		"movieId":   req.MovieId,
		"userId":    req.UserId,
		"text":      req.Text,
		"createdAt": req.CreatedAt.AsTime(),
	}

	id, err := e.Comments.AddCommentService(commentMap)
	if err != nil {
		return &pb.AddCommentResponse{Err: err.Error()}, err
	}

	m := metrics.New()
	m.NewComments()

	return &pb.AddCommentResponse{Id: id}, nil
}

func (e *Endpoint) UpdateComment(_ context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id комментария не указан")
	} else if req.Text == "" {
		return nil, errors.New("текст комментария не указан")
	}

	commentMap := map[string]interface{}{
		"id":        req.Id,
		"text":      req.Text,
		"updatedAt": req.UpdatedAt.AsTime(),
	}

	err := e.Comments.UpdateCommentService(commentMap)
	if err != nil {
		return &pb.UpdateCommentResponse{Err: err.Error()}, err
	}

	return &pb.UpdateCommentResponse{Err: ""}, nil
}

func (e *Endpoint) DelComment(_ context.Context, req *pb.DelCommentRequest) (*pb.DelCommentResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("id новости не указан")
	}

	err := e.Comments.DelCommentService(req.Id)
	if err != nil {
		return &pb.DelCommentResponse{Err: err.Error()}, err
	}

	return &pb.DelCommentResponse{Err: ""}, nil
}

func convertCommentsToPb(comments []models.Comments) []*pb.GetCommentsByIdItem {
	var pbComments []*pb.GetCommentsByIdItem
	for _, comment := range comments {
		pbComment := &pb.GetCommentsByIdItem{
			Id:       comment.ID,
			ParentId: comment.ParentID,
			User: &pb.GetUserById{
				Id:        comment.User.ID,
				Username:  comment.User.Username,
				PhotoUrl:  comment.User.PhotoUrl,
				FirstName: comment.User.FirstName,
				LastName:  comment.User.LastName,
			},
			Text:      comment.Text,
			CreatedAt: ts.New(comment.CreatedAt),
			UpdatedAt: ts.New(comment.UpdatedAt),
			Children:  convertCommentsToPb(comment.Children),
		}
		pbComments = append(pbComments, pbComment)
	}
	return pbComments
}
