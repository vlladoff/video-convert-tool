package repository

import (
	"context"
	"github.com/vlladoff/video-convert-tool/internal/storage"
	"os"
)

type VideoRepository struct {
	storage storage.Storage
}

func NewVideoRepository(storage storage.Storage) *VideoRepository {
	return &VideoRepository{storage: storage}
}

func (vr *VideoRepository) SaveVideo(ctx context.Context, path, s3Path string) error {
	videoData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return vr.storage.Upload(ctx, s3Path, videoData)
}
