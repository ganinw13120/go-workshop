// go:build integration
package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositorySuite struct {
	suite.Suite

	pool     *dockertest.Pool
	resource *dockertest.Resource

	mongoClient    *mongo.Client
	userCollection *mongo.Collection

	userRepository repository.IUser
}

func (s *UserRepositorySuite) SetupSuite() {
	var err error

	s.pool, err = dockertest.NewPool("")

	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	if err = s.pool.Client.Ping(); err != nil {
		s.FailNow("Could not ping docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	s.resource, err = s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		s.FailNow("Could not start resource: %s", err)
	}

	err = s.pool.Retry(func() error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var err error
		s.mongoClient, err = mongo.Connect(
			ctx,
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", s.resource.GetPort("27017/tcp")),
			),
		)

		if err != nil {
			return err
		}

		return s.mongoClient.Ping(ctx, nil)
	})

	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}
}

func (s *UserRepositorySuite) TearDownSuite() {
	if err := s.pool.Purge(s.resource); err != nil {
		s.FailNow("Could not purge resource: %s", err)
	}

	if err := s.mongoClient.Disconnect(context.Background()); err != nil {
		s.FailNow("Could not disconnect mongo client: %s", err)
	}
}

func (s *UserRepositorySuite) SetupTest() {
	s.userCollection = s.mongoClient.Database("test").Collection("users")
	mongoAdapter := adapter.NewMongoDBAdapter(s.mongoClient)
	s.userRepository = repository.NewUser(
		repository.UserConfig{
			Timeout: 10 * time.Second,
		},
		mongoAdapter,
		s.userCollection,
	)
}

func (s *UserRepositorySuite) TearDownTest() {
	s.userCollection.Drop(context.Background())
}

func (s *UserRepositorySuite) TestGetAll() {

	s.Run("should return empty slice when collection is empty", func() {
		users, err := s.userRepository.GetAll()

		s.NoError(err)
		s.Empty(users)
	})

	s.Run("should return all users", func() {
		// insert 2 users
		_, err := s.userCollection.InsertMany(context.Background(), []interface{}{
			map[string]interface{}{
				"name": "user1",
			},
			map[string]interface{}{
				"name": "user2",
			},
		})

		s.NoError(err)

		users, err := s.userRepository.GetAll()

		s.NoError(err)
		s.Len(users, 2)
		s.Equal("user1", users[0].Name)
		s.Equal("user2", users[1].Name)
	})
}

func (s *UserRepositorySuite) TestGetByID() {

	s.Run("should return error when id is invalid", func() {
		_, err := s.userRepository.GetByID("")
		s.Error(err)
	})

	s.Run("should return error when user not found", func() {
		_, err := s.userRepository.GetByID("63dccac268616ec85ccfcfd2")

		s.Error(err)
	})

	s.Run("should return user when user found", func() {
		obj, _ := s.userCollection.InsertOne(context.Background(), map[string]interface{}{
			"name": "user1",
		})

		user, err := s.userRepository.GetByID(obj.InsertedID.(primitive.ObjectID).Hex())

		s.NoError(err)
		s.Equal("user1", user.Name)
	})
}

func (s *UserRepositorySuite) TestCreate() {
	s.Run("should return user", func() {
		user := entity.User{
			ID:        "1",
			Name:      "user1",
			Username:  "user1",
			Password:  "password",
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		id, err := s.userRepository.Create(&user)

		s.NoError(err)
		s.NotEmpty(id)
	})
}

func (s *UserRepositorySuite) TestUpdate() {
	s.Run("should return error when id is invalid", func() {
		_, err := s.userRepository.Update("", &entity.User{})

		s.Error(err)
	})

	s.Run("should return error when user not found", func() {
		_, err := s.userRepository.Update("63dccac268616ec85ccfcfd2", &entity.User{})

		s.Error(err)
	})

	s.Run("should return user when user found", func() {
		obj, _ := s.userCollection.InsertOne(context.Background(), map[string]interface{}{
			"name": "user1",
		})

		user := entity.User{
			ID:   obj.InsertedID.(primitive.ObjectID).Hex(),
			Name: "user2",
		}

		isSuccess, err := s.userRepository.Update(user.ID, &user)

		s.Run("should return true when user updated", func() {
			s.NoError(err)
			s.Equal(true, isSuccess)

		})

		s.Run("should return user with updated name", func() {
			updateUser := entity.User{}
			err = s.userCollection.FindOne(context.Background(), bson.M{"_id": obj.InsertedID}).Decode(&updateUser)

			s.NoError(err)
			s.Equal("user2", user.Name)
		})
	})
}

func (s *UserRepositorySuite) TestDelete() {
	s.Run("should return error when id is invalid", func() {
		err := s.userRepository.Delete("")

		s.Error(err)
	})

	s.Run("should return error when user not found", func() {
		err := s.userRepository.Delete("63dccac268616ec85ccfcfd2")

		s.Error(err)
	})

	s.Run("should return true when user found", func() {
		obj, _ := s.userCollection.InsertOne(context.Background(), map[string]interface{}{
			"name": "user1",
		})

		err := s.userRepository.Delete(obj.InsertedID.(primitive.ObjectID).Hex())

		s.Run("no error should be returned", func() {
			s.NoError(err)
		})

		s.Run("user should be deleted", func() {
			user := entity.User{}
			err = s.userCollection.FindOne(context.Background(), bson.M{"_id": obj.InsertedID}).Decode(&user)

			s.Error(err)
		})
	})
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
