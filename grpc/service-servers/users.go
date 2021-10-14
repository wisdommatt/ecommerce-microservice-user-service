package servers

import (
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/mappers"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/panick"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	userService users.Service
}

// NewUserService returns a new user service.
func NewUserService(userService users.Service) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// CreateUser is the grpc handler to create new user.
func (u *UserServiceServer) CreateUser(ctx context.Context, req *proto.NewUser) (res *proto.User, err error) {
	globalTracer := opentracing.GlobalTracer()
	span := globalTracer.StartSpan("create-user")
	defer span.Finish()
	defer panick.RecoverFromPanic(opentracing.ContextWithSpan(ctx, span))
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("time", time.Now())
	span.LogFields(log.Object("request.body", req))

	ctx = opentracing.ContextWithSpan(ctx, span)
	newUser, err := u.userService.CreateUser(ctx, mappers.ProtoNewUserToInternalUser(req))
	if err != nil {
		return nil, errors.New("An error occured, please try again later !")
	}
	return mappers.InternalToProtoUser(newUser), nil
}
