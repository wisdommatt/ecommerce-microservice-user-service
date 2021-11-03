package servers

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
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
	span := opentracing.StartSpan("CreateUser")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("request.body", req)

	ctx = opentracing.ContextWithSpan(ctx, span)
	newUser, err := u.userService.CreateUser(ctx, ProtoNewUserToInternalUser(req))
	if err != nil {
		return nil, err
	}
	return InternalToProtoUser(newUser), nil
}

func (u *UserServiceServer) GetUsers(ctx context.Context, filter *proto.GetUsersFilter) (*proto.GetUsersResponse, error) {
	span := opentracing.StartSpan("GetUsers")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("param.filter", filter)

	ctx = opentracing.ContextWithSpan(ctx, span)
	users, err := u.userService.GetUsers(ctx, filter.AfterId, filter.Limit)
	if err != nil {
		return nil, err
	}
	var protoUsers []*proto.User
	for _, user := range users {
		protoUsers = append(protoUsers, InternalToProtoUser(&user))
	}
	return &proto.GetUsersResponse{
		Users: protoUsers,
	}, nil
}

func (u *UserServiceServer) LoginUser(ctx context.Context, input *proto.LoginInput) (*proto.LoginResponse, error) {
	span := opentracing.StartSpan("LoginUser")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("param.input", input)

	ctx = opentracing.ContextWithSpan(ctx, span)
	usr, jwtToken, err := u.userService.LoginUser(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{
		User:     InternalToProtoUser(usr),
		JwtToken: jwtToken,
	}, nil
}

func (u *UserServiceServer) GetUserFromJWT(ctx context.Context, input *proto.GetUserFromJWTInput) (*proto.GetUserFromJWTResponse, error) {
	span := opentracing.StartSpan("GetUserFromJWT")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("param.input", input)

	ctx = opentracing.ContextWithSpan(ctx, span)
	usr, err := u.userService.GetUserFromJWT(ctx, input.JwtToken)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserFromJWTResponse{
		User: InternalToProtoUser(usr),
	}, nil
}
