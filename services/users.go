package services

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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
	LoginUser(ctx context.Context, email, password string) (*users.User, string, error)
	GetUserFromJWT(ctx context.Context, jwtToken string) (*users.User, error)
}

type UserServiceImpl struct {
	userRepo users.Repository
	natsConn *nats.Conn
	tracer   opentracing.Tracer
}

var (
	ErrPaginationLimit = errors.New("pagination limit max is 100")
	ErrTryAgain        = errors.New("an error occured, please try again later")
)

// NewUserService returns a new user service.
func NewUserService(userRepo users.Repository, tracer opentracing.Tracer, natsConn *nats.Conn) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
		natsConn: natsConn,
		tracer:   tracer,
	}
}

// CreateUser is the service handler to create new user.
func (s *UserServiceImpl) CreateUser(ctx context.Context, newUser *users.User) (*users.User, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateUser")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
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
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetUsers")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
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

func (s *UserServiceImpl) LoginUser(ctx context.Context, email, password string) (*users.User, string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LoginUser")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	if email == "" || password == "" {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("error.object", "some field are empty"),
			log.Event("input validation"),
		)
		return nil, "", errors.New("all fields are required")
	}
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	if user == nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("error.object", "user with email does not exist"))
		return nil, "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(err),
			log.Event("password validation"),
		)
		return nil, "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.ID,
		"timeAdded": user.TimeAdded,
		"exp":       time.Now().AddDate(0, 0, 4).UTC().UnixNano(),
	})
	jwtToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("jwt generation"))
		return nil, "", ErrTryAgain
	}
	return user, jwtToken, nil
}

func (s *UserServiceImpl) GetUserFromJWT(ctx context.Context, jwtToken string) (*users.User, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetUserFromJWT")
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid jwt token string")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("jwt decoding"))
		return nil, errors.New("invalid jwt")
	}
	if !token.Valid {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("invalid jwt token")),
			log.Event("jwt token validation"),
		)
		return nil, errors.New("jwt token is not valid")
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)
	span.SetTag("jwtClaims", claims)
	user, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	return user, nil
}
