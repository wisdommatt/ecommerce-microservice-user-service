package services

import (
	"context"
	"errors"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	userMocks "github.com/wisdommatt/ecommerce-microservice-user-service/test/mocks/users"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepoMock := &userMocks.RepositoryMock{}
	tests := []struct {
		name               string
		newUser            *users.User
		want               *users.User
		repoCreateUserFunc func(ctx context.Context, user *users.User) error
		wantErr            bool
	}{
		{
			name: "CreateUser repository implementation with error",
			newUser: &users.User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
			},
			repoCreateUserFunc: func(ctx context.Context, user *users.User) error {
				return errors.New("Invalid user entity !")
			},
			wantErr: true,
		},
		{
			name: "CreateUser repository implementation without error",
			newUser: &users.User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
			},
			repoCreateUserFunc: func(ctx context.Context, user *users.User) error {
				user.ID = "john.doe"
				return nil
			},
			want: &users.User{
				ID:       "john.doe",
				FullName: "John Doe",
				Country:  "Nigeria",
			},
		},
		{
			name: "empty password",
			newUser: &users.User{
				FullName: "John Doe",
				Password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock.CreateUserFunc = tt.repoCreateUserFunc
			s := NewUserService(userRepoMock)
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
