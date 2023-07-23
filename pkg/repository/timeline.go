package repository

import (
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
)

type ITimeline interface {
	GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error)
	Save(entity.Thread) error
}

type timeline struct {
	mongoDBAdapter adapter.IMongoDBAdapter
}

func NewTimeline(mongoDBAdapter adapter.IMongoDBAdapter) *timeline {
	return &timeline{
		mongoDBAdapter: mongoDBAdapter,
	}
}

func (t timeline) GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error) {
	return nil, nil
}

func (t timeline) Save(thread entity.Thread) error {
	return nil
}
