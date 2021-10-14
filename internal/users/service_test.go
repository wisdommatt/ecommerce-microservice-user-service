package users

import (
	"context"
	"errors"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepoMock := &RepositoryMock{}
	tests := []struct {
		name               string
		newUser            *User
		want               *User
		repoCreateUserFunc func(ctx context.Context, user *User) error
		wantErr            bool
	}{
		{
			name: "CreateUser repository implementation with error",
			newUser: &User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
			},
			repoCreateUserFunc: func(ctx context.Context, user *User) error {
				return errors.New("Invalid user entity !")
			},
			wantErr: true,
		},
		{
			name: "CreateUser repository implementation without error",
			newUser: &User{
				FullName: "John Doe",
				Country:  "Nigeria",
				Password: "123456",
			},
			repoCreateUserFunc: func(ctx context.Context, user *User) error {
				user.ID = "john.doe"
				return nil
			},
			want: &User{
				ID:       "john.doe",
				FullName: "John Doe",
				Country:  "Nigeria",
			},
		},
		{
			name: "empty password",
			newUser: &User{
				FullName: "John Doe",
				Password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock.CreateUserFunc = tt.repoCreateUserFunc
			s := NewService(userRepoMock)
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
