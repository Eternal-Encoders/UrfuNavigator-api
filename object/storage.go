package object

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
)

func (s *MinIOS3) putObject(key string, fileData []byte, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	reader := bytes.NewReader(fileData)
	opts := minio.PutObjectOptions{ContentType: contentType}

	_, err := s.Client.PutObject(ctx, s.BucketName, key, reader, reader.Size(), opts)
	return err
}

func (s *MinIOS3) removeObject(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	return s.Client.RemoveObject(ctx, s.BucketName, key, minio.RemoveObjectOptions{})
}

func (s *MinIOS3) objectExists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	_, err := s.Client.StatObject(ctx, s.BucketName, key, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}

	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return false, nil
	}

	return false, err
}
