package object

import (
	"net"
	"net/http"
	"time"
	"urfunavigator/index/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ConnectOptions struct {
	RequestTimeout time.Duration
	PresignTTL     time.Duration
	MaxRetries     int
}

type MinIOS3 struct {
	Client         *minio.Client
	BucketName     string
	requestTimeout time.Duration
	presignTTL     time.Duration
	maxRetries     int
}

func Connect(endpoint string, access string, secret string, bucketName string, opts ConnectOptions) *MinIOS3 {
	if opts.RequestTimeout <= 0 {
		opts.RequestTimeout = 60 * time.Second
	}
	if opts.PresignTTL <= 0 {
		opts.PresignTTL = time.Hour
	}
	if opts.MaxRetries < 0 {
		opts.MaxRetries = 0
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 20
	transport.IdleConnTimeout = 90 * time.Second
	transport.ResponseHeaderTimeout = opts.RequestTimeout
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	option := &minio.Options{
		Creds:     credentials.NewStaticV4(access, secret, ""),
		Secure:    true,
		Region:    "us-east-1",
		Transport: transport,
	}

	minioClient, err := minio.New(endpoint, option)
	if err != nil {
		logger.Fatal("failed to connect to minio", "endpoint", endpoint, "bucket", bucketName, "err", err)
	}

	logger.Info("connected to minio",
		"endpoint", endpoint,
		"bucket", bucketName,
		"request_timeout_sec", opts.RequestTimeout.Seconds(),
		"presign_ttl_sec", opts.PresignTTL.Seconds(),
		"max_retries", opts.MaxRetries,
	)

	return &MinIOS3{
		Client:         minioClient,
		BucketName:     bucketName,
		requestTimeout: opts.RequestTimeout,
		presignTTL:     opts.PresignTTL,
		maxRetries:     opts.MaxRetries,
	}
}

func (s *MinIOS3) Disconnect() {}
