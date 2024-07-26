package object

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOS3 struct {
	Client     *minio.Client
	BucketName string
}

func Connect(endpoint string, access string, secret string, bucketName string) *MinIOS3 {
	option := &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: true,
	}
	minioClient, err := minio.New(endpoint, option)

	if err != nil {
		log.Fatal(err)
	}

	return &MinIOS3{
		Client:     minioClient,
		BucketName: bucketName,
	}
}

func (s *MinIOS3) GetFile(fileName string) ([]byte, error) {
	option := minio.GetObjectOptions{}

	file, err := s.Client.GetObject(context.TODO(), s.BucketName, fileName, option)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 4)
	res := []byte{}
	for {
		_, err := file.Read(buf)
		res = append(res, buf...)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return res, nil
		}
	}
}
