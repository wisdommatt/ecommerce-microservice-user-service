package services

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser *users.User) (*users.User, error)
	GetUsers(ctx context.Context, afterId string, limit int32) ([]users.User, error)
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

func (s *UserServiceImpl) GetUsers(ctx context.Context, afterId string, limit int32) ([]users.User, error) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = opentracing.StartSpan("service.GetUsers")
	}
	if limit == 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("event", "no filter limit provided"),
		)
		return nil, errors.New("filter limit must be provided")
	}
	if limit > 100 {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errPaginationLimit),
		)
		return nil, errPaginationLimit
	}
	users, err := s.userRepo.GetUsers(ctx, afterId, limit)
	if err != nil {
		return nil, errors.New("an error occured, please try again later")
	}
	return users, nil
}
