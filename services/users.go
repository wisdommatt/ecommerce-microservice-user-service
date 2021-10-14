package services

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser *users.User) (*users.User, error)
}

type UserServiceImpl struct {
	userRepo users.Repository
}

// NewUserService returns a new user service.
func NewUserService(userRepo users.Repository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser is the service handler to create new user.
func (s *UserServiceImpl) CreateUser(ctx context.Context, newUser *users.User) (*users.User, error) {
	var err error
	newUser.Password, err = password.HashPassword(ctx, newUser.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
