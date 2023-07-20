package usecase

import (
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type IUser interface {
	GetAll() ([]entity.User, error)
	GetByID(id string) (entity.User, error)
	Create(user *entity.User) (string, error)
	Update(id string, user *entity.User) (bool, error)
	Delete(id string) error
}

type user struct {
	repo repository.IUser
}

func NewUser(repo repository.IUser) IUser {
	return &user{
		repo,
	}
}

func (u user) GetAll() ([]entity.User, error) {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u user) GetByID(id string) (entity.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u user) Create(user *entity.User) (string, error) {
	userID, err := u.repo.Create(user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (u user) Update(id string, user *entity.User) (bool, error) {
	isSuccess, err := u.repo.Update(id, user)
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}

func (u user) Delete(id string) error {
	return u.repo.Delete(id)
}
