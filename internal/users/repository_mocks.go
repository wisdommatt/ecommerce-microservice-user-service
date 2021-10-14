package users

import (
	"context"
)

type RepositoryMock struct {
	CreateUserFunc func(ctx context.Context, user *User) error
}

func (r *RepositoryMock) CreateUser(ctx context.Context, user *User) error {
	return r.CreateUserFunc(ctx, user)
}
