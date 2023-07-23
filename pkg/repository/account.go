package repository

import (
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
)

type IAccount interface {
	GetAccount() ([]*entity.Account, error)
}

type account struct {
	mongoDBAdapter adapter.IMongoDBAdapter
}

func NewAccount(mongoDBAdapter adapter.IMongoDBAdapter) *account {
	return &account{
		mongoDBAdapter: mongoDBAdapter,
	}
}

func (t timeline) GetAccount() ([]*entity.Account, error) {
	return nil, nil
}
