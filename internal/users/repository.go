package users

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/conversions"
	"github.com/wisdommatt/ecommerce-microservice-user-service/pkg/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context, afterId string, limit int32) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
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

	tracer.SetMongoDBSpanComponentTags(span, r.collection.Name())
	span.SetTag("param.newUser", conversions.ToJSON(span, newUser))
	_, err := r.collection.InsertOne(ctx, newUser)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("error.object", err.Error(), "event", "mongodb.InsertOne")
		return err
	}
	return nil
}

func (r *UserRepo) GetUsers(ctx context.Context, afterId string, limit int32) ([]User, error) {
	span := r.tracer.StartSpan("GetUsers", opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()
	tracer.SetMongoDBSpanComponentTags(span, r.collection.Name())

	filter := bson.M{
		"_id": bson.M{"$gt": afterId},
	}
	findOpts := options.Find().SetLimit(int64(limit))
	span.SetTag("param.afterId", afterId).SetTag("param.limit", limit)
	span.SetTag("mongodb.filter", conversions.ToJSON(span, filter))

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("error.object", err.Error(), "event", "mongodb.Find")
		return nil, err
	}
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("error.object", err.Error(), "event", "mongodb.Cursor.All")
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	span := r.tracer.StartSpan("GetUserByEmail", opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()
	tracer.SetMongoDBSpanComponentTags(span, r.collection.Name())

	filter := bson.M{"email": email}
	span.SetTag("param.email", email).SetTag("mongodb.filter", conversions.ToJSON(span, filter))
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("mongodb.FindOne"))
		return nil, err
	}
	return &user, nil
}
