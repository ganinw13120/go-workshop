package usecase

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type IThread interface {
	Save(context.Context, entity.Thread) error
	GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error)
}

type thread struct {
	threadRepo repository.IThread
}

func NewThread(threadRepo repository.IThread) *thread {
	return &thread{
		threadRepo: threadRepo,
	}
}

func (u thread) GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error) {
	return u.threadRepo.GetTimelineFromHashtag(ctx, hashtag, cursor, pageSize)
}

func (u thread) Save(ctx context.Context, thread entity.Thread) error {
	return u.threadRepo.Save(ctx, thread)
}
