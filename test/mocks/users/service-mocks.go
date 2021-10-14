package users

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

type ServiceMock struct {
	CreateUserFunc func(ctx context.Context, newUser *users.User) (*users.User, error)
}

func (s *ServiceMock) CreateUser(ctx context.Context, newUser *users.User) (*users.User, error) {
	return s.CreateUserFunc(ctx, newUser)
}
