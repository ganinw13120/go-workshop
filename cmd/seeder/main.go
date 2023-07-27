package main

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"

	"github.com/wisesight/go-api-template/config"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/repository"
	"github.com/wisesight/go-api-template/pkg/validator"
)

func main() {

	keywords := []string{
		"alpskub",
		"tiger",
		"tangtai",
		"barbie",
		"alpsbarbie",
		"cutegirl",
		"helloworld",
		"sadboy",
	}

	cfg := config.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbClient, err := adapter.NewMongoDBConnection(ctx, cfg.MongoDBURI)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err = validator.NewValidator(); err != nil {
		panic(err)
	}

	accountCollection := mongodbClient.Database("go-workshop").Collection("accounts")
	threadCollection := mongodbClient.Database("go-workshop").Collection("threads")
	mongoDBAdapter := adapter.NewMongoDBAdapter(mongodbClient)

	threadRepo := repository.NewThread(mongoDBAdapter, threadCollection)
	accountRepo := repository.NewAccount(mongoDBAdapter, accountCollection)

	threadUseCase := usecase.NewThread(cfg, threadRepo)
	accountUseCase := usecase.NewAccount(accountRepo)

	_ = threadUseCase
	_ = accountUseCase
	accounts := make([]entity.Account, 50)
	for i := 0; i < 50; i++ {
		account := entity.Account{
			Id:              primitive.NewObjectID(),
			DisplayName:     faker.Name(),
			Username:        faker.Username(),
			ProfileImageUrl: faker.URL(),
			Description:     faker.Sentence(),
			Follower:        rand.Intn(1000),
			Following:       rand.Intn(1000),
		}
		fmt.Println(account)
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := accountUseCase.Save(ctx, account)
		if err != nil {
			panic(err)
		}
		accounts[i] = account
	}

	var threads []entity.Thread
	for i := 0; i < 1000; i++ {
		var parentThread *primitive.ObjectID
		keyword := keywords[rand.Intn(len(keywords))]
		if len(threads) > 0 && rand.Intn(10) == 0 {
			parent := threads[rand.Intn(len(threads))].Id
			parentThread = &parent
		}
		ind := rand.Intn(len(accounts))
		userId := accounts[ind].Id
		thread := entity.Thread{
			Id:           primitive.NewObjectID(),
			Text:         fmt.Sprintf("%s #%s", faker.Sentence(), keyword),
			UserId:       userId,
			Likes:        rand.Intn(10000),
			ParentThread: parentThread,
		}
		fmt.Println(thread)
		if len(threads) < 10 {
			threads = append(threads, thread)
		} else if rand.Intn(3) == 0 {
			threads = append(threads, thread)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := threadUseCase.Save(ctx, thread)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(accounts)
}
