package repository

import (
	"context"
	"time"

	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/apperror"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUser interface {
	GetAll() ([]entity.User, error)
	GetByID(id string) (entity.User, error)
	Create(user *entity.User) (string, error)
	Update(id string, user *entity.User) (bool, error)
	Delete(id string) error
}

type UserConfig struct {
	Timeout time.Duration
}

type user struct {
	mongoDBAdapter adapter.IMongoDBAdapter
	userCollection adapter.IMongoCollection
	timeout        time.Duration
}

func NewUser(userConfig UserConfig, mongoDBAdapter adapter.IMongoDBAdapter, userCollection adapter.IMongoCollection) IUser {

	return &user{
		mongoDBAdapter: mongoDBAdapter,
		userCollection: userCollection,
		timeout:        userConfig.Timeout,
	}
}

func (r user) GetAll() ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var users []entity.User

	err := r.mongoDBAdapter.Find(ctx, r.userCollection, &users, bson.D{})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r user) GetByID(id string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	primitiveObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entity.User{}, err
	}

	var user entity.User

	err = r.mongoDBAdapter.FindOne(ctx, r.userCollection, &user, bson.D{{Key: "_id", Value: primitiveObjectID}})

	if err != nil {
		// if err == adapter.ErrNoDocuments {
		// return entity.User{}, ErrUserNotFound
		// }
		return entity.User{}, err
	}

	return user, nil
}

func (r user) Create(user *entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// validate

	primitiveObjectID, err := r.mongoDBAdapter.InsertOne(ctx, r.userCollection, user)

	if err != nil {
		return "", err
	}

	return primitiveObjectID.Hex(), nil
}

func (r user) Update(id string, user *entity.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	primitiveObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	isSuccess, err := r.mongoDBAdapter.UpdateOne(ctx, r.userCollection, bson.D{{Key: "_id", Value: primitiveObjectID}}, bson.D{{Key: "$set", Value: user}})

	if err != nil {
		return false, err
	}

	if !isSuccess {
		return false, apperror.NewError(
			"User not found",
			"User not found",
			apperror.NotFound,
		)
	}

	return true, nil
}

func (r user) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	primitiveObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	isSuccess, err := r.mongoDBAdapter.DeleteOne(ctx, r.userCollection, bson.D{{Key: "_id", Value: primitiveObjectID}})

	if err != nil {
		return err
	}

	if !isSuccess {
		return apperror.NewError(
			"User not found",
			"User not found",
			apperror.NotFound,
		)
	}

	return nil
}
