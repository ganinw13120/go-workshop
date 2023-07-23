package repository

import (
	"context"
	"fmt"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IThread interface {
	GetTimelineFromHashtag(context.Context, string, *string, int) ([]entity.Thread, *string, error)
	Save(context.Context, entity.Thread) error
}

type thread struct {
	mongoDBAdapter   adapter.IMongoDBAdapter
	threadCollection adapter.IMongoCollection
}

func NewThread(mongoDBAdapter adapter.IMongoDBAdapter, threadCollection adapter.IMongoCollection) *thread {
	return &thread{
		mongoDBAdapter:   mongoDBAdapter,
		threadCollection: threadCollection,
	}
}

func (t thread) GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", 1}})
	opt.SetLimit(int64(pageSize))
	var filter bson.M
	if cursor != nil {
		lastCursor, err := primitive.ObjectIDFromHex(*cursor)
		if err != nil {
			return nil, nil, err
		}
		filter = bson.M{"_id": bson.M{"$gte": lastCursor}, "text": bson.M{"$regex": fmt.Sprintf("#%s", hashtag)}}
	} else {
		filter = bson.M{"text": bson.M{"$regex": fmt.Sprintf("#%s", hashtag)}}
	}
	var threads []entity.Thread
	err := t.mongoDBAdapter.Find(ctx, t.threadCollection, &threads, filter, opt)
	if err != nil {
		return nil, nil, err
	}
	var nextPage *string
	if len(threads) >= pageSize {
		var nextThread entity.Thread
		nextPageFilter := bson.M{"_id": bson.M{"$gt": threads[pageSize-1].Id}, "text": bson.M{"$regex": fmt.Sprintf("#%s", hashtag)}}
		nextOpt := options.FindOne()
		opt.SetSort(bson.D{{"_id", 1}})
		err = t.mongoDBAdapter.FindOne(ctx, t.threadCollection, &nextThread, nextPageFilter, nextOpt)
		if err == nil {
			nextPageId := nextThread.Id.Hex()
			nextPage = &nextPageId
		}
	}
	return threads, nextPage, nil
}

func (t thread) Save(ctx context.Context, thread entity.Thread) error {
	thread.Id = primitive.NewObjectID()
	_, err := t.mongoDBAdapter.InsertOne(ctx, t.threadCollection, thread)
	return err
}
