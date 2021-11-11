package servers

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/mocks"
)

func TestUserServiceServer_CreateUser(t *testing.T) {
	userService := &mocks.UserService{}

	userService.On("CreateUser", mock.Anything, ProtoNewUserToInternalUser(&proto.NewUser{
		FullName: "John Doe",
		Email:    "john.doe@example.com",
		Password: "123456",
		Country:  "Nigeria",
	})).Return(nil, errors.New("an erorr occured"))

	userService.On("CreateUser", mock.Anything, ProtoNewUserToInternalUser(&proto.NewUser{
		FullName: "Jane Doe",
		Email:    "jane.doe@example.com",
		Password: "123456",
		Country:  "Nigeria",
	})).Return(&users.User{
		ID:       "jane.doe123",
		FullName: "Jane Doe",
		Email:    "jane.doe@example.com",
		Country:  "Nigeria",
	}, nil)

	tests := []struct {
		name                  string
		req                   *proto.NewUser
		serviceCreateUserFunc func(ctx context.Context, newUser *users.User) (*users.User, error)
		wantRes               *proto.User
		wantErr               bool
	}{
		{
			name: "CreateUser service implementation that returns an error",
			req: &proto.NewUser{
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Password: "123456",
				Country:  "Nigeria",
			},
			wantErr: true,
		},
		{
			name: "CreateUser service implementation without error",
			req: &proto.NewUser{
				FullName: "Jane Doe",
				Email:    "jane.doe@example.com",
				Password: "123456",
				Country:  "Nigeria",
			},
			wantRes: &proto.User{
				Id:       "jane.doe123",
				FullName: "Jane Doe",
				Email:    "jane.doe@example.com",
				Country:  "Nigeria",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserServiceServer(userService)
			gotRes, err := u.CreateUser(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceServer.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("UserServiceServer.CreateUser() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestUserServiceServer_GetUsers(t *testing.T) {
	userService := &mocks.UserService{}
	userService.On("GetUsers", mock.Anything, "invalid", int32(100)).Return(nil, errors.New("an error occured"))
	userService.On("GetUsers", mock.Anything, "valid", int32(3)).Return([]users.User{
		{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
	}, nil)
	userService.On("GetUsers", mock.Anything, "empty", int32(0)).Return(nil, nil)

	tests := []struct {
		name                string
		filter              *proto.GetUsersFilter
		serviceGetUsersFunc func(ctx context.Context, afterId string, limit int32) ([]users.User, error)
		want                *proto.GetUsersResponse
		wantErr             bool
	}{
		{
			name:    "GetUsers service implementation with error",
			filter:  &proto.GetUsersFilter{AfterId: "invalid", Limit: 100},
			wantErr: true,
		},
		{
			name:   "GetUsers service implementation without error",
			filter: &proto.GetUsersFilter{AfterId: "valid", Limit: 3},
			want: &proto.GetUsersResponse{
				Users: []*proto.User{
					{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
				},
			},
		},
		{
			name:   "GetUsers service implementation with empty reponse",
			filter: &proto.GetUsersFilter{AfterId: "empty"},
			want:   &proto.GetUsersResponse{Users: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserServiceServer(userService)
			got, err := u.GetUsers(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceServer.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceServer.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_LoginUser(t *testing.T) {
	userService := &mocks.UserService{}
	userService.On("LoginUser", mock.Anything, "invalid@example.com", "123456").
		Return(nil, "", errors.New("an error occured"))
	userService.On("LoginUser", mock.Anything, "valid@example.com", "123456").
		Return(&users.User{
			ID:       "valid.user",
			FullName: "Valid User",
		}, "theJwtToken", nil)

	type args struct {
		input *proto.LoginInput
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.LoginResponse
		wantErr bool
	}{
		{
			name: "LoginUser service implementation with error",
			args: args{input: &proto.LoginInput{
				Email:    "invalid@example.com",
				Password: "123456",
			}},
			wantErr: true,
		},
		{
			name: "LoginUser service implementation without error",
			args: args{input: &proto.LoginInput{
				Email:    "valid@example.com",
				Password: "123456",
			}},
			want: &proto.LoginResponse{
				User: &proto.User{
					Id:       "valid.user",
					FullName: "Valid User",
				},
				JwtToken: "theJwtToken",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserServiceServer(userService)
			got, err := u.LoginUser(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceServer.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceServer.LoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_GetUserFromJWT(t *testing.T) {
	userService := &mocks.UserService{}
	userService.On("GetUserFromJWT", mock.Anything, "invalidJwtToken").Return(nil, errors.New("an error occured"))
	userService.On("GetUserFromJWT", mock.Anything, "validJwtToken").Return(&users.User{
		ID:       "valid.user",
		FullName: "Valid User",
		Email:    "valid@example.com",
	}, nil)

	type args struct {
		input *proto.GetUserFromJWTInput
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.GetUserFromJWTResponse
		wantErr bool
	}{
		{
			name: "GetUserFromJWT service implementation with error",
			args: args{input: &proto.GetUserFromJWTInput{
				JwtToken: "invalidJwtToken",
			}},
			wantErr: true,
		},
		{
			name: "GetUserFromJWT service implementation with error",
			args: args{input: &proto.GetUserFromJWTInput{
				JwtToken: "validJwtToken",
			}},
			want: &proto.GetUserFromJWTResponse{
				User: &proto.User{
					Id:       "valid.user",
					FullName: "Valid User",
					Email:    "valid@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserServiceServer(userService)
			got, err := u.GetUserFromJWT(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceServer.GetUserFromJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceServer.GetUserFromJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
