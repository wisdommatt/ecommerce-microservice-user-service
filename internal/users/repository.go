package users

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/panick"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
}

type UserRepo struct {
	collection *mongo.Collection
	tracer     opentracing.Tracer
}

// NewRepository returns a new user repository object that implements the
// Repository interface.
func NewRepository(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
		tracer:     tracer.Init("mongodb"),
	}
}

// CreateUser adds a new user to the database.
func (r *UserRepo) CreateUser(ctx context.Context, newUser *User) error {
	newUser.ID = primitive.NewObjectID().Hex()
	newUser.TimeAdded = time.Now()
	newUser.LastUpdated = time.Now()
	span := r.tracer.StartSpan("CreateUser", opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()
	defer panick.RecoverFromPanic(opentracing.ContextWithSpan(ctx, span))

	tracer.SetMongoDBSpanComponentTags(span, r.collection.Name())
	span.SetTag("param.newUser", newUser)
	span.LogKV("event", "about to save the new user to the database", "time", time.Now())
	_, err := r.collection.InsertOne(ctx, newUser)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(err),
			log.Event("mongodb.InsertOne"),
		)
		return err
	}
	span.LogKV("event", "new user saved successfully", "time", time.Now())
	return nil
}
