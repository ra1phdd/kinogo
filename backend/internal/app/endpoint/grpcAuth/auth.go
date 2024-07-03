package grpcAuth

import (
	"context"
	pb "kinogo/pkg/auth_v1"
)

type Auth interface {
	ValidateTelegramAuth(data map[string]interface{}, botToken string) bool
	AddUserIfNotExists(data map[string]interface{})
	GenerateToken(data map[string]interface{}, jwtSecret string) (string, error)
	ValidateToken(tokenString string, data map[string]interface{}, jwtSecret string) (bool, error)
	CheckAdminService(id int32) (bool, error)
}

type Endpoint struct {
	Auth      Auth
	JwtSecret string
	pb.UnimplementedAuthV1Server
}

func (e *Endpoint) CheckAuth(_ context.Context, req *pb.CheckAuthRequest) (*pb.CheckAuthResponse, error) {
	authMap := map[string]interface{}{
		"id":         req.UserId,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"username":   req.Username,
		"photo_url":  req.PhotoUrl,
		"auth_date":  req.AuthDate,
		"isAdmin":    req.IsAdmin,
	}

	isAuth, err := e.Auth.ValidateToken(req.Token, authMap, e.JwtSecret)
	if err != nil {
		return &pb.CheckAuthResponse{}, err
	}

	return &pb.CheckAuthResponse{
		IsAuth: isAuth,
		Err:    "",
	}, nil
}
