package server

import (
	"context"

	"github.com/Skapar/backend/internal/service"
	stock "github.com/Skapar/backend/proto"
)

type Server struct {
	stock.UnsafeStockServiceServer

	Service service.Service
}

// New Server constructor.
func New(
	service service.Service,
) *Server {
	return &Server{
		Service: service,
	}
}

func (s *Server) CreateUser(ctx context.Context, in *stock.CreateUserRequest) (*stock.CreateUserResponse, error) {
	userID, err := s.Service.CreateUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return &stock.CreateUserResponse{
		UserId: userID,
	}, nil
}
