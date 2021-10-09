package users

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
}

type UserRepo struct {
	collection *mongo.Collection
}

// NewRepository returns a new user repository object that implements the
// Repository interface.
func NewRepository(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

// Create adds a new user to the database.
func (r *UserRepo) Create(ctx context.Context, user *User) error {
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(string)
	return nil
}
