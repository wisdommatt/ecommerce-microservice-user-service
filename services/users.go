package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
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

var (
	ErrPaginationLimit = errors.New("pagination limit max is 100")
	ErrTryAgain        = errors.New("an error occured, please try again later")
)

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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	span.SetTag("bcrypt.passwordCost", bcrypt.DefaultCost)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("param.passwordStr", newUser.Password)
		span.LogFields(log.Error(err), log.String("event", "password hash error"))
		return nil, ErrTryAgain
	}
	newUser.Password = string(passwordHash)
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
	natsMessageJSON, err := json.Marshal(natsMessage)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("converting object to json"), log.Object("object", natsMessage))
		return
	}
	span.SetTag("nats.message", string(natsMessageJSON))
	err = s.natsConn.Publish("notification.SendEmail", natsMessageJSON)
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
		span.LogFields(log.Error(ErrPaginationLimit))
		return nil, ErrPaginationLimit
	}
	users, err := s.userRepo.GetUsers(ctx, afterId, limit)
	if err != nil {
		return nil, ErrTryAgain
	}
	return users, nil
}
