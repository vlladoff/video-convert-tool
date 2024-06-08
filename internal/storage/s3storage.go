package storage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vlladoff/video-convert-tool/internal/config"
	"io"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(cfg config.Config) (*MinIOStorage, error) {
	client, err := minio.New(cfg.MinioS3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioS3.AccessKey, cfg.MinioS3.SecretKey, ""),
		Secure: cfg.MinioS3.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return &MinIOStorage{
		client: client,
		bucket: cfg.MinioS3.Bucket,
	}, nil
}

func (m *MinIOStorage) Upload(ctx context.Context, key string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := m.client.PutObject(ctx, m.bucket, key, reader, reader.Size(), minio.PutObjectOptions{})
	return err
}

func (m *MinIOStorage) Download(ctx context.Context, key string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		// todo log
		return nil, err
	}

	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			// todo log
		}
	}(object)

	return io.ReadAll(object)
}

func (m *MinIOStorage) Delete(ctx context.Context, key string) error {
	return m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
}
