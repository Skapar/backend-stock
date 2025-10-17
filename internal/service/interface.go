package service

import (
	"context"

	"github.com/Skapar/backend/internal/models/entities"
	pb "github.com/Skapar/backend/proto"
)

type Service interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]*entities.User, error)
}
