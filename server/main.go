package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wisdommatt/ecommerce-microservice-user-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-user-service/internal/users"
	"github.com/wisdommatt/ecommerce-microservice-user-service/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	serviceTracer, closer := initJaeger("user-service", log)
	defer closer.Close()
	opentracing.SetGlobalTracer(serviceTracer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2020"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.WithError(err).Fatal("TCP conn error")
	}
	mongoDBClient := mustConnectMongoDB()
	userRepository := users.NewRepository(mongoDBClient)

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, services.NewUserService(userRepository))
	log.Info("Server running on port: ", port)
	grpcServer.Serve(lis)
}

func mustConnectMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(os.Getenv("MONGODB_DATABASE_NAME"))
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}

func initJaeger(service string, log *logrus.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.WithError(err).Fatal("ERROR: cannot init Jaeger")
	}
	return tracer, closer
}
