package usecase

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type IAccount interface {
	Get(ctx context.Context, id string) (*entity.Account, error)
	Save(context.Context, entity.Account) error
}

type account struct {
	accountRepo repository.IAccount
}

func NewAccount(accountRepo repository.IAccount) *account {
	return &account{
		accountRepo: accountRepo,
	}
}

func (u account) Save(ctx context.Context, account entity.Account) error {
	return u.accountRepo.Save(ctx, account)
}

func (u account) Get(ctx context.Context, id string) (*entity.Account, error) {
	result, err := u.accountRepo.GetAccount(ctx, id)
	return result, err
}
