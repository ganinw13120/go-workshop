package usecase

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type ITimeline interface {
	Save(context.Context, entity.Thread) error
	GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error)
}

type timeline struct {
	timelineRepo repository.ITimeline
}

func NewTimeline(timelineRepo repository.ITimeline) *timeline {
	return &timeline{
		timelineRepo: timelineRepo,
	}
}

func (u timeline) GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error) {
	return u.timelineRepo.GetTimelineFromHashtag(ctx, hashtag, cursor, pageSize)
}

func (u timeline) Save(ctx context.Context, thread entity.Thread) error {
	return u.timelineRepo.Save(ctx, thread)
}
