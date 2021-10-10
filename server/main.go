package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/panick"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/tracer"
	"github.com/wisdommatt/ecommerce-microservice-user-service/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	serviceTracer := tracer.Init("user-service")
	opentracing.SetGlobalTracer(serviceTracer)
	panicSpan := serviceTracer.StartSpan("user-service-panic")
	defer panicSpan.Finish()
	defer panick.RecoverFromPanic(opentracing.ContextWithSpan(context.Background(), panicSpan))

	port := os.Getenv("PORT")
	if port == "" {
		port = "2020"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.WithError(err).Fatal("TCP conn error")
	}
	mongoDBClient := mustConnectMongoDB(log)
	userRepository := users.NewRepository(mongoDBClient)
	userService := users.NewService(userRepository)

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, services.NewUserService(userService))
	log.Info("Server running on port: ", port)
	grpcServer.Serve(lis)
}

func mustConnectMongoDB(log *logrus.Logger) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to mongodb")
	}
	return client.Database(os.Getenv("MONGODB_DATABASE_NAME"))
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}
