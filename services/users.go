package services

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	userRepo users.Repository
}

// NewUserService returns a new user service.
func NewUserService(userRepo users.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser is the rpc handler create
func (u *UserService) CreateUser(ctx context.Context, req *proto.NewUser) (res *proto.User, err error) {
	return res, err
}
