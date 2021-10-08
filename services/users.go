package services

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
}

// NewUserService returns a new user service.
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser is the rpc handler create
func (u *UserService) CreateUser(ctx context.Context, req *proto.NewUser) (res *proto.User, err error) {
	return res, err
}
