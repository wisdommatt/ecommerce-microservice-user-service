package services

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/conversions"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser *users.User) (*users.User, error)
	GetUsers(ctx context.Context, afterId string, limit int32) ([]users.User, error)
}

type UserServiceImpl struct {
	userRepo users.Repository
	natsConn *nats.Conn
}

// NewUserService returns a new user service.
func NewUserService(userRepo users.Repository, natsConn *nats.Conn) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
		natsConn: natsConn,
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
	s.publishCreateUserSendEmailEvent(span, newUser)
	return newUser, nil
}

func (s *UserServiceImpl) publishCreateUserSendEmailEvent(span opentracing.Span, user *users.User) {
	span = opentracing.StartSpan("publish-create-user-email-event", opentracing.ChildOf(span.Context()))
	defer span.Finish()
	natsMessage := map[string]string{
		"to":      user.Email,
		"subject": "Welcome to my microservice application",
		"body":    "It's glad to have you onboard, thanks for checking it out",
	}
	span.SetTag("nats.message", conversions.ToJSON(span, natsMessage))
	err := s.natsConn.Publish("notification.SendEmail", []byte(conversions.ToJSON(span, natsMessage)))
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("nats.notification.SendEmail"))
	}
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
