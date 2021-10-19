package servers

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/service-servers/mappers"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/panick"
	"github.com/wisdommatt/ecommerce-microservice-user-service/services"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	userService services.UserService
}

// NewUserServiceServer returns a new user service.
func NewUserServiceServer(userService services.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// CreateUser is the grpc handler to create new user.
func (u *UserServiceServer) CreateUser(ctx context.Context, req *proto.NewUser) (res *proto.User, err error) {
	span := opentracing.GlobalTracer().StartSpan("CreateUser")
	defer span.Finish()
	defer panick.RecoverFromPanic(opentracing.ContextWithSpan(ctx, span))
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("time", time.Now())
	span.LogFields(log.Object("request.body", req))

	ctx = opentracing.ContextWithSpan(ctx, span)
	newUser, err := u.userService.CreateUser(ctx, mappers.ProtoNewUserToInternalUser(req))
	if err != nil {
		return nil, err
	}
	return mappers.InternalToProtoUser(newUser), nil
}

func (u *UserServiceServer) GetUsers(ctx context.Context, filter *proto.GetUsersFilter) (*proto.GetUsersResponse, error) {
	span := opentracing.GlobalTracer().StartSpan("GetUsers")
	defer span.Finish()
	defer panick.RecoverFromPanic(opentracing.ContextWithSpan(ctx, span))
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("time", time.Now())
	span.SetTag("param.filter", filter)

	ctx = opentracing.ContextWithSpan(ctx, span)
	users, err := u.userService.GetUsers(ctx, filter.AfterId, filter.Limit)
	if err != nil {
		return nil, err
	}
	var protoUsers []*proto.User
	for _, user := range users {
		protoUsers = append(protoUsers, mappers.InternalToProtoUser(&user))
	}
	return &proto.GetUsersResponse{
		Users: protoUsers,
	}, nil
}
