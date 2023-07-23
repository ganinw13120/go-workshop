package usecase

import (
	"context"
	"github.com/wisesight/go-api-template/config"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/helper"
	"github.com/wisesight/go-api-template/pkg/repository"
)

type IThread interface {
	Save(context.Context, entity.Thread) error
	GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error)
}

type thread struct {
	config     config.Config
	threadRepo repository.IThread
}

func NewThread(config config.Config, threadRepo repository.IThread) *thread {
	return &thread{
		config:     config,
		threadRepo: threadRepo,
	}
}

func (u thread) GetTimelineFromHashtag(ctx context.Context, hashtag string, cursor *string, pageSize int) ([]entity.Thread, *string, error) {
	if cursor != nil {
		decryptedCursor, err := helper.Decrypt(u.config.EncryptedKey, *cursor)
		if err != nil {
			return nil, nil, err
		}
		cursor = &decryptedCursor
	}
	result, nextCursor, err := u.threadRepo.GetTimelineFromHashtag(ctx, hashtag, cursor, pageSize)
	if err != nil {
		return nil, nil, err
	}
	if nextCursor != nil {
		encryptedCursor, err := helper.Encrypt(u.config.EncryptedKey, *nextCursor)
		if err != nil {
			return nil, nil, err
		}
		nextCursor = &encryptedCursor
	}
	return result, nextCursor, nil
}

func (u thread) Save(ctx context.Context, thread entity.Thread) error {
	return u.threadRepo.Save(ctx, thread)
}
