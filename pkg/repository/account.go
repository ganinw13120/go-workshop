package repository

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccount interface {
	GetAccount() ([]*entity.Account, error)
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

func (t account) GetAccount() ([]*entity.Account, error) {
	return nil, nil
}

func (t account) Save(ctx context.Context, account entity.Account) error {
	account.Id = primitive.NewObjectID()
	_, err := t.mongoDBAdapter.InsertOne(ctx, t.accountCollection, account)
	return err
}
