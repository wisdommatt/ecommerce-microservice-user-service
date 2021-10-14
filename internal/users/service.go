package users

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, newUser *User) (*User, error)
}

type UserService struct {
	userRepo Repository
}

// NewService returns a new user service.
func NewService(userRepo Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser is the service handler to create new user.
func (s *UserService) CreateUser(ctx context.Context, newUser *User) (*User, error) {
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
