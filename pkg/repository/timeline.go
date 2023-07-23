package repository

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITimeline interface {
	GetTimelineFromHashtag(context.Context, string, *string, int) ([]entity.Thread, *string, error)
	Save(context.Context, entity.Thread) error
}

type timeline struct {
	mongoDBAdapter   adapter.IMongoDBAdapter
	threadCollection adapter.IMongoCollection
}

func NewTimeline(mongoDBAdapter adapter.IMongoDBAdapter, threadCollection adapter.IMongoCollection) *timeline {
	return &timeline{
		mongoDBAdapter:   mongoDBAdapter,
		threadCollection: threadCollection,
	}
}

func (t timeline) GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", 1}})
	opt.SetLimit(int64(pageSize))
	var filter bson.M
	if cursor != nil {
		lastCursor, err := primitive.ObjectIDFromHex(*cursor)
		if err != nil {
			return nil, nil, err
		}
		filter = bson.M{"_id": bson.M{"$gte": lastCursor}}
	}
	var threads []entity.Thread
	err := t.mongoDBAdapter.Find(ctx, t.threadCollection, &threads, filter, opt)
	if err != nil {
		return nil, nil, err
	}
	var nextPage *string
	if len(threads) >= pageSize {
		var nextThread entity.Thread
		nextPageFilter := bson.M{"_id": bson.M{"$gt": threads[pageSize-1].Id}}
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

func (t timeline) Save(ctx context.Context, thread entity.Thread) error {
	thread.Id = primitive.NewObjectID()
	_, err := t.mongoDBAdapter.InsertOne(ctx, t.threadCollection, thread)
	return err
}
