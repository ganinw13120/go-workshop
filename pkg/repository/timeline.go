package repository

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITimeline interface {
	GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error)
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

func (t timeline) GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error) {
	return nil, nil
}

func (t timeline) Save(ctx context.Context, thread entity.Thread) error {
	thread.Id = primitive.NewObjectID()
	_, err := t.mongoDBAdapter.InsertOne(ctx, t.threadCollection, thread)
	return err
}
