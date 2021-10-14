package servers

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	userMocks "github.com/wisdommatt/ecommerce-microservice-user-service/test/mocks/users"
)

func TestUserServiceServer_CreateUser(t *testing.T) {
	userServiceMock := &userMocks.ServiceMock{}
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
			serviceCreateUserFunc: func(ctx context.Context, newUser *users.User) (*users.User, error) {
				return nil, errors.New("Password is too weak !")
			},
			wantErr: true,
		},
		{
			name: "CreateUser service implementation without error",
			req: &proto.NewUser{
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Password: "123456",
				Country:  "Nigeria",
			},
			serviceCreateUserFunc: func(ctx context.Context, newUser *users.User) (*users.User, error) {
				return &users.User{
					ID:       "john.doe123",
					FullName: "John Doe",
					Email:    "john.doe@example.com",
					Password: "123456",
					Country:  "Nigeria",
				}, nil
			},
			wantRes: &proto.User{
				Id:       "john.doe123",
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Country:  "Nigeria",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock.CreateUserFunc = tt.serviceCreateUserFunc
			u := NewUserService(userServiceMock)
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
