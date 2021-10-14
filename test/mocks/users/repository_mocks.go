package users

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

type RepositoryMock struct {
	CreateUserFunc func(ctx context.Context, user *users.User) error
}

func (r *RepositoryMock) CreateUser(ctx context.Context, user *users.User) error {
	return r.CreateUserFunc(ctx, user)
}
