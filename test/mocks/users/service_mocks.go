package users

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

type ServiceMock struct {
	CreateUserFunc func(ctx context.Context, newUser *users.User) (*users.User, error)
	GetUsersFunc   func(ctx context.Context, afterId string, limit int32) ([]users.User, error)
	LoginUserFunc  func(ctx context.Context, email, password string) (*users.User, string, error)
}

func (s *ServiceMock) CreateUser(ctx context.Context, newUser *users.User) (*users.User, error) {
	return s.CreateUserFunc(ctx, newUser)
}

func (s *ServiceMock) GetUsers(ctx context.Context, afterId string, limit int32) ([]users.User, error) {
	return s.GetUsersFunc(ctx, afterId, limit)
}

func (s *ServiceMock) LoginUser(ctx context.Context, email, password string) (*users.User, string, error) {
	return s.LoginUserFunc(ctx, email, password)
}
