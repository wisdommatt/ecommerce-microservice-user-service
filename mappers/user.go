package mappers

import (
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
)

func InternalToProtoUser(usr *users.User) *proto.User {
	return &proto.User{
		Id:       usr.ID,
		FullName: usr.FullName,
		Email:    usr.Email,
		Country:  usr.Country,
	}
}

func ProtoNewUserToInternalUser(usr *proto.NewUser) *users.User {
	return &users.User{
		FullName: usr.FullName,
		Email:    usr.Email,
		Country:  usr.Country,
	}
}
