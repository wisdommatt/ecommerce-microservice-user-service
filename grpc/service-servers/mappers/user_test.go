package mappers

import (
	"reflect"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

func TestInternalToProtoUser(t *testing.T) {
	type args struct {
		usr *users.User
	}
	tests := []struct {
		name string
		args args
		want *proto.User
	}{
		{
			name: "all fields",
			args: args{
				usr: &users.User{
					ID:       "hello",
					FullName: "Wisdom Matt",
					Email:    "hello@example.com",
					Country:  "Nigeria",
				},
			},
			want: &proto.User{
				Id:       "hello",
				FullName: "Wisdom Matt",
				Email:    "hello@example.com",
				Country:  "Nigeria",
			},
		},
		{
			name: "incomplete fields",
			args: args{
				usr: &users.User{
					ID:       "hhh",
					FullName: "Wisdom Matt",
				},
			},
			want: &proto.User{
				Id:       "hhh",
				FullName: "Wisdom Matt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InternalToProtoUser(tt.args.usr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InternalToProtoUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtoNewUserToInternalUser(t *testing.T) {
	type args struct {
		usr *proto.NewUser
	}
	tests := []struct {
		name string
		args args
		want *users.User
	}{
		{
			name: "all fields",
			args: args{
				usr: &proto.NewUser{
					FullName: "Wisdom Matt",
					Email:    "hello@example.com",
					Country:  "Nigeria",
					Password: "123456",
				},
			},
			want: &users.User{
				FullName: "Wisdom Matt",
				Email:    "hello@example.com",
				Country:  "Nigeria",
				Password: "123456",
			},
		},
		{
			name: "incomplete fields",
			args: args{
				usr: &proto.NewUser{
					FullName: "Wisdom Matt",
					Email:    "hello@example.com",
				},
			},
			want: &users.User{
				FullName: "Wisdom Matt",
				Email:    "hello@example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoNewUserToInternalUser(tt.args.usr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoNewUserToInternalUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
