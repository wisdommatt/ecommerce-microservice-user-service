package users

import (
	"context"
)

type ServiceMock struct {
	CreateUserFunc func(ctx context.Context, newUser *User) (*User, error)
}

func (s *ServiceMock) CreateUser(ctx context.Context, newUser *User) (*User, error) {
	return s.CreateUserFunc(ctx, newUser)
}
