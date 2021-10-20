package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	servers "github.com/wisdommatt/ecommerce-microservice-user-service/grpc/service-servers"
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

	natsConn, err := nats.Connect(os.Getenv("NATS_URI"))
	if err != nil {
		log.WithField("nats_uri", os.Getenv("NATS_URI")).WithError(err).
			Error("an error occured while connecting to nats")
	}
	defer natsConn.Close()

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
	userService := services.NewUserService(userRepository, natsConn)

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, servers.NewUserServiceServer(userService))
	log.WithField("nats_uri", os.Getenv("NATS_URI")).Info("Server running on port: ", port)
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
