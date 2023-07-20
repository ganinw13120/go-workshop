package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository/mocks"
	"github.com/wisesight/go-api-template/pkg/usecase"
)

type UserUsecaseSuite struct {
	suite.Suite
	userRepo    *mocks.IUser
	userUseCase usecase.IUser

	resUserRepoGetAll []entity.User
	errUserRepoGetAll error

	resUserRepoGetByID entity.User
	errUserRepoGetByID error

	resUserRepoCreate string
	errUserRepoCreate error

	resUserRepoUpdate bool
	errUserRepoUpdate error

	errUserRepoDelete error
}

func (s *UserUsecaseSuite) SetupSuite() {
	s.userRepo = &mocks.IUser{}
	s.userUseCase = usecase.NewUser(s.userRepo)

	s.userRepo.On("GetAll").Return(
		func() []entity.User {
			return s.resUserRepoGetAll
		},
		func() error {
			return s.errUserRepoGetAll
		},
	)

	s.userRepo.On("GetByID", mock.Anything).Return(
		func(string) entity.User {
			return s.resUserRepoGetByID
		},
		func(string) error {
			return s.errUserRepoGetByID
		},
	)

	s.userRepo.On("Create", mock.Anything).Return(
		func(*entity.User) string {
			return s.resUserRepoCreate
		},
		func(*entity.User) error {
			return s.errUserRepoCreate
		},
	)

	s.userRepo.On("Update", mock.Anything, mock.Anything).Return(
		func(string, *entity.User) bool {
			return s.resUserRepoUpdate
		},
		func(string, *entity.User) error {
			return s.errUserRepoUpdate
		},
	)

	s.userRepo.On("Delete", mock.Anything).Return(
		func(string) error {
			return s.errUserRepoDelete
		},
	)
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

func (s *UserUsecaseSuite) SetupTest() {
	s.resUserRepoGetAll = []entity.User{}
	s.errUserRepoGetAll = nil

	s.resUserRepoGetByID = entity.User{}
	s.errUserRepoGetByID = nil

	s.resUserRepoCreate = "id"
	s.errUserRepoCreate = nil

	s.resUserRepoUpdate = true
	s.errUserRepoUpdate = nil

	s.errUserRepoDelete = nil
}

func (s *UserUsecaseSuite) TestGetAll() {

	s.Run("should get all user", func() {
		s.userUseCase.GetAll()
		s.userRepo.AssertCalled(s.T(), "GetAll")
	})

	s.Run("should return error when user failed", func() {
		s.resUserRepoGetAll = []entity.User{}
		s.errUserRepoGetAll = errors.New("get all failed")

		res, err := s.userUseCase.GetAll()

		s.Nil(res)
		s.EqualError(err, "get all failed")
	})

	s.Run("should return user when get all success", func() {
		s.resUserRepoGetAll = []entity.User{
			{
				ID:   "mock-id",
				Name: "test",
			},
		}
		s.errUserRepoGetAll = nil

		res, err := s.userUseCase.GetAll()

		s.Equal(res, s.resUserRepoGetAll)
		s.Nil(err)
	})

}

func (s *UserUsecaseSuite) TestGetByID() {

	s.Run("should get user by id", func() {
		id := "mock-id"

		s.userUseCase.GetByID(id)

		s.userRepo.AssertCalled(s.T(), "GetByID", id)
	})

	s.Run("should return error when get user by id failed", func() {
		id := "mock-id"
		s.resUserRepoGetByID = entity.User{}
		s.errUserRepoGetByID = errors.New("get by id failed")

		res, err := s.userUseCase.GetByID(id)

		s.Equal(res, entity.User{})
		s.EqualError(err, "get by id failed")
	})

	s.Run("should return user when get user by id success", func() {
		id := "mock-id"
		s.resUserRepoGetByID = entity.User{
			ID:   id,
			Name: "test",
		}
		s.errUserRepoGetByID = nil

		res, err := s.userUseCase.GetByID(id)

		s.Equal(res, s.resUserRepoGetByID)
		s.Nil(err)
	})
}

func (s *UserUsecaseSuite) TestCreate() {

	s.Run("should create user", func() {
		user := entity.User{
			Name: "test",
		}

		s.userUseCase.Create(&user)

		s.userRepo.AssertCalled(s.T(), "Create", &user)
	})

	s.Run("should return error when create user failed", func() {
		user := entity.User{
			Name: "test",
		}
		s.resUserRepoCreate = ""
		s.errUserRepoCreate = errors.New("create failed")

		_, err := s.userUseCase.Create(&user)

		s.EqualError(err, "create failed")
	})

	s.Run("should return user when create user success", func() {
		user := entity.User{
			ID:   "mock-id",
			Name: "test",
		}
		s.resUserRepoCreate = "mock-id"
		s.errUserRepoCreate = nil

		res, err := s.userUseCase.Create(&user)

		s.Equal(res, "mock-id")
		s.Nil(err)
	})
}

func (s *UserUsecaseSuite) TestUpdate() {

	s.Run("should update user", func() {
		id := "mock-id"
		user := entity.User{
			Name: "test",
		}

		s.userUseCase.Update(id, &user)

		s.userRepo.AssertCalled(s.T(), "Update", id, &user)
	})

	s.Run("should return error when update user failed", func() {
		id := "mock-id"
		user := entity.User{
			Name: "test",
		}

		s.resUserRepoUpdate = false
		s.errUserRepoUpdate = errors.New("update failed")

		res, err := s.userUseCase.Update(id, &user)

		s.Equal(res, false)
		s.EqualError(err, "update failed")
	})

	s.Run("should return user when update user success", func() {
		id := "mock-id"
		user := entity.User{
			ID:   id,
			Name: "test",
		}
		s.resUserRepoUpdate = true
		s.errUserRepoUpdate = nil

		res, err := s.userUseCase.Update(id, &user)

		s.Equal(res, true)
		s.Nil(err)
	})
}

func (s *UserUsecaseSuite) TestDelete() {

	s.Run("should delete user", func() {
		id := "mock-id"

		s.userUseCase.Delete(id)

		s.userRepo.AssertCalled(s.T(), "Delete", id)
	})

	s.Run("should return error when delete user failed", func() {
		id := "mock-id"
		s.errUserRepoDelete = errors.New("delete failed")

		err := s.userUseCase.Delete(id)

		s.EqualError(err, "delete failed")
	})

	s.Run("should return nil when delete user success", func() {
		id := "mock-id"
		s.errUserRepoDelete = nil

		err := s.userUseCase.Delete(id)

		s.Nil(err)
	})
}
