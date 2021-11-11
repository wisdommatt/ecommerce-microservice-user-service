package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/mocks"
	userMocks "github.com/wisdommatt/ecommerce-microservice-user-service/test/mocks/users"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepo := &mocks.Repository{}
	userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*users.User")).Once().Return(nil).Run(func(args mock.Arguments) {
		usr := args[1].(*users.User)
		usr.Password = "hashedPassword"
		usr.ID = "john.doe"
	})
	userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*users.User")).Return(errors.New("invalid user"))
	userRepo.On("GetUserByEmail", mock.Anything, "valid@example.com").Return(nil, nil)
	userRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(&users.User{FullName: "User"}, nil)
	userRepo.On("GetUserByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("an error occured"))

	tests := []struct {
		name    string
		newUser *users.User
		want    *users.User
		wantErr bool
	}{
		{
			name: "CreateUser repository implementation without error",
			newUser: &users.User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
				Email:    "valid@example.com",
			},
			want: &users.User{
				ID:       "john.doe",
				FullName: "John Doe",
				Country:  "Nigeria",
				Email:    "valid@example.com",
				Password: "hashedPassword",
			},
		},
		{
			name: "CreateUser repository implementation with error",
			newUser: &users.User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
				Email:    "valid@example.com",
			},
			wantErr: true,
		},
		{
			name: "existing email",
			newUser: &users.User{
				Email: "existing@example.com",
			},
			wantErr: true,
		},
		{
			name: "GetUserByEmail repository implementation with error",
			newUser: &users.User{
				Email: "error@example.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(userRepo, nil)
			got, err := s.CreateUser(context.Background(), tt.newUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && (tt.want.ID != got.ID || tt.want.FullName != got.FullName || tt.want.Email != got.Email) {
				t.Errorf("UserServiceServer.CreateUser() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestUserServiceImpl_GetUsers(t *testing.T) {
	userRepoMock := &userMocks.RepositoryMock{}
	type args struct {
		afterId string
		limit   int32
	}
	tests := []struct {
		name             string
		args             args
		repoGetUsersFunc func(ctx context.Context, afterId string, limit int32) ([]users.User, error)
		want             []users.User
		wantErr          bool
	}{
		{
			name:    "no pagination limit",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "pagination > 100",
			args:    args{limit: 101},
			wantErr: true,
		},
		{
			name: "GetUsers repo implementation with error",
			args: args{limit: 100},
			repoGetUsersFunc: func(ctx context.Context, afterId string, limit int32) ([]users.User, error) {
				return nil, errors.New("DB connection error !")
			},
			wantErr: true,
		},
		{
			name: "testcase with no expected error",
			args: args{limit: 10},
			repoGetUsersFunc: func(ctx context.Context, afterId string, limit int32) ([]users.User, error) {
				return []users.User{
					{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
				}, nil
			},
			want: []users.User{
				{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock.GetUsersFunc = tt.repoGetUsersFunc
			s := NewUserService(userRepoMock, nil)
			got, err := s.GetUsers(context.Background(), tt.args.afterId, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("User, nilServiceImpl.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceImpl.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
