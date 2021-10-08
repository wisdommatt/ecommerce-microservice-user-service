package main

import (
	"log"
	"net"
	"os"

	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/services"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2020"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("TCP conn error: ", err.Error())
	}
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, services.NewUserService())
	log.Println("Server running on port: ", port)
	grpcServer.Serve(lis)
}
