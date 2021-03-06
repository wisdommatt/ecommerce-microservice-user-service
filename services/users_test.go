package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/mocks"
	"golang.org/x/crypto/bcrypt"
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
			s := NewUserService(userRepo, &opentracing.NoopTracer{}, nil)
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
	userRepo := &mocks.Repository{}
	userRepo.On("GetUsers", mock.Anything, "", int32(100)).Return(nil, errors.New("an error occured"))
	userRepo.On("GetUsers", mock.Anything, "valid", int32(2)).Return([]users.User{
		{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
	}, nil)

	type args struct {
		afterId string
		limit   int32
	}
	tests := []struct {
		name    string
		args    args
		want    []users.User
		wantErr bool
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
			name:    "GetUsers repo implementation with error",
			args:    args{limit: 100},
			wantErr: true,
		},
		{
			name: "testcase with no expected error",
			args: args{limit: 2, afterId: "valid"},
			want: []users.User{
				{FullName: "John"}, {FullName: "Jane"}, {FullName: "Doe"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(userRepo, &opentracing.NoopTracer{}, nil)
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

func TestUserServiceImpl_LoginUser(t *testing.T) {
	userHashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)

	userRepo := &mocks.Repository{}
	userRepo.On("GetUserByEmail", mock.Anything, "invalid@example.com").Return(nil, errors.New("an error occured"))
	userRepo.On("GetUserByEmail", mock.Anything, "nil@example.com").Return(nil, nil)
	userRepo.On("GetUserByEmail", mock.Anything, "valid@example.com").Return(&users.User{
		ID:       "valid.user",
		FullName: "Valid User",
		Password: string(userHashedPassword),
	}, nil)

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			name:    "empty email",
			args:    args{email: "", password: "123456"},
			wantErr: true,
		},
		{
			name:    "empty password",
			args:    args{email: "user@example.com", password: ""},
			wantErr: true,
		},
		{
			name:    "empty email and password",
			args:    args{email: "", password: ""},
			wantErr: true,
		},
		{
			name:    "GetUserByEmail repo implementation with error",
			args:    args{email: "invalid@example.com", password: "123456"},
			wantErr: true,
		},
		{
			name:    "GetUserByEmail repo implementation with nil user response",
			args:    args{email: "nil@example.com", password: "123456"},
			wantErr: true,
		},
		{
			name:    "invalid password",
			args:    args{email: "valid@example.com", password: "1234567"},
			wantErr: true,
		},
		{
			name: "valid credentials",
			args: args{email: "valid@example.com", password: "123456"},
			want: &users.User{
				ID:       "valid.user",
				FullName: "Valid User",
				Password: string(userHashedPassword),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(userRepo, &opentracing.NoopTracer{}, nil)
			got, got1, err := s.LoginUser(context.Background(), tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImpl.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceImpl.LoginUser() got = %v, want %v", got, tt.want)
			}
			if !tt.wantErr && got1 == "" {
				t.Errorf("UserServiceImpl.LoginUser() jwtToken = %v should not be empty", got1)
			}
		})
	}
}

func TestUserServiceImpl_GetUserFromJWT(t *testing.T) {
	userRepo := &mocks.Repository{}
	userRepo.On("GetUserByID", mock.Anything, "user.invalid").Return(nil, errors.New("an error occured"))
	userRepo.On("GetUserByID", mock.Anything, "user.valid").Return(&users.User{
		ID:       "user.valid",
		FullName: "Valid User",
	}, nil)

	type args struct {
		jwtToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			name: "invalid jwt token",
			args: args{
				jwtToken: "invalidJwtToken",
			},
			wantErr: true,
		},
		{
			name: "expired jwt token",
			args: args{
				jwtToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIxLTExLTA3VDEzOjM3OjA1LjM2OTgyODQ1KzAxOjAwIiwidXNlcklkIjoidXNlci5pbnZhbGlkIn0.R538RkECQiyoja9SqU37xYg78m1a0ONJgp2i-g_ZLMU",
			},
			wantErr: true,
		},
		{
			name: "GetUserByID repo implementation with error",
			args: args{
				jwtToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyLmludmFsaWQifQ.EJjgUmsONbvp9W-lfdmrSNsUy402LWEOgFgNiEobiXc",
			},
			wantErr: true,
		},
		{
			name: "GetUserByID repo implementation without error",
			args: args{
				jwtToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyLnZhbGlkIn0.s7PBgz9LPIpolroL8bH7NfR1u1b3UC8ra0cmooOPQ9I",
			},
			want: &users.User{
				ID:       "user.valid",
				FullName: "Valid User",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(userRepo, &opentracing.NoopTracer{}, nil)
			got, err := s.GetUserFromJWT(context.Background(), tt.args.jwtToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImpl.GetUserFromJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceImpl.GetUserFromJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
