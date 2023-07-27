package repository

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccount interface {
	GetAccount(ctx context.Context, id string) (*entity.Account, error)
	Save(context.Context, entity.Account) error
}

type account struct {
	mongoDBAdapter    adapter.IMongoDBAdapter
	accountCollection adapter.IMongoCollection
}

func NewAccount(mongoDBAdapter adapter.IMongoDBAdapter, accountCollection adapter.IMongoCollection) *account {
	return &account{
		mongoDBAdapter:    mongoDBAdapter,
		accountCollection: accountCollection,
	}
}

func (t account) GetAccount(ctx context.Context, id string) (*entity.Account, error) {
	var account entity.Account
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = t.mongoDBAdapter.FindOne(ctx, t.accountCollection, &account, bson.M{"_id": userId}, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (t account) Save(ctx context.Context, account entity.Account) error {
	_, err := t.mongoDBAdapter.InsertOne(ctx, t.accountCollection, account)
	return err
}
