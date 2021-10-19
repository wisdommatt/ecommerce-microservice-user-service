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
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = opentracing.StartSpan("service.GetUsers")
	}
	existingUser, err := s.userRepo.GetUserByEmail(ctx, newUser.Email)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("existing user email validation"))
		return nil, ErrTryAgain
	}
	if existingUser != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("error.object", "user with email already exist"),
			log.Event("existing user email validation"),
		)
		return nil, errors.New("user with this email already exist")
	}
	newUser.Password, err = password.HashPassword(ctx, newUser.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrTryAgain
	}
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, ErrTryAgain
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
			log.Error(ErrPaginationLimit),
		)
		return nil, ErrPaginationLimit
	}
	users, err := s.userRepo.GetUsers(ctx, afterId, limit)
	if err != nil {
		return nil, ErrTryAgain
	}
	return users, nil
}
