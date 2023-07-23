package usecase

import (
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type ITimeline interface {
	Save(entity.Thread) error
	GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error)
}

type timeline struct {
	timelineRepo repository.ITimeline
}

func NewTimeline(timelineRepo repository.ITimeline) *timeline {
	return &timeline{
		timelineRepo: timelineRepo,
	}
}

func (t timeline) GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error) {
	return t.GetTimelineFromHashtag(hashtag)
}

func (t timeline) Save(thread entity.Thread) error {
	return t.Save(thread)
}
