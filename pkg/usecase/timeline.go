package usecase

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type ITimeline interface {
	Save(context.Context, entity.Thread) error
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

func (u timeline) GetTimelineFromHashtag(hashtag string) ([]*entity.Thread, error) {
	return u.timelineRepo.GetTimelineFromHashtag(hashtag)
}

func (u timeline) Save(ctx context.Context, thread entity.Thread) error {
	return u.timelineRepo.Save(ctx, thread)
}
