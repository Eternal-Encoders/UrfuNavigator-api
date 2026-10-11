package config

import (
	"os"
	"strconv"
	"strings"
	"time"
	"urfunavigator/index/logger"
)

type Config struct {
	Mode         string
	Port         string
	Cors         []string
	DefaultPath  string
	LogLevel     string
	LogFormat    string
	DbUri        string
	DbCollection string
	S3Endpoint           string
	S3Access             string
	S3Secret             string
	S3bucketName         string
	S3RequestTimeout     time.Duration
	S3PresignTTL         time.Duration
	S3MaxRetries         int
	JWTSecret            string
	JWTExpiresIn         time.Duration
}

func New() *Config {
	return &Config{
		Mode:         getEnvOptional("MODE", "DEV"),
		Port:         getEnv("PORT"),
		Cors:         getEnvArrayStr("CORS"),
		DefaultPath:  getEnv("DEFAULT_PATH"),
		LogLevel:     getEnvOptional("LOG_LEVEL", "info"),
		LogFormat:    getEnvOptional("LOG_FORMAT", "text"),
		DbUri:        getEnv("DATABASE_URI"),
		DbCollection: getEnv("DATABASE_COLLECTION"),
		S3Endpoint:       getEnv("BUCKET_ENDPOINT"),
		S3Access:         getEnv("BUCKET_ACCESS_KEY"),
		S3Secret:         getEnv("BUCKET_SECRET_KEY"),
		S3bucketName:     getEnv("BUCKET_NAME"),
		S3RequestTimeout: time.Duration(getEnvIntOptional("S3_REQUEST_TIMEOUT_SEC", 60)) * time.Second,
		S3PresignTTL:     time.Duration(getEnvIntOptional("S3_PRESIGN_TTL_SEC", 3600)) * time.Second,
		S3MaxRetries:     getEnvIntOptional("S3_MAX_RETRIES", 2),
		JWTSecret:        getEnv("JWT_SECRET"),
		JWTExpiresIn:     time.Duration(getEnvIntOptional("JWT_EXPIRES_IN_HOURS", 24)) * time.Hour,
	}
}

func getEnvOptional(key string, fallback string) string {
	value, exist := os.LookupEnv(key)
	if !exist || strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func getEnv(key string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		file, fileExist := os.LookupEnv(key + "_FILE")

		if !fileExist {
			logger.Fatal("required environment variable is missing", "key", key)
		}

		data, err := os.ReadFile(file)
		if err != nil {
			logger.Fatal("failed to read environment file", "key", key, "err", err)
		}
		return string(data)
	}

	return value
}

func getEnvInt(key string) int {
	value := getEnv(key)

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		logger.Fatal("environment variable is not an integer", "key", key, "err", err)
	}

	return valueInt
}

func getEnvIntOptional(key string, fallback int) int {
	value, exist := os.LookupEnv(key)
	if !exist || strings.TrimSpace(value) == "" {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		logger.Warn("invalid integer env value, using fallback", "key", key, "fallback", fallback, "err", err)
		return fallback
	}

	return valueInt
}

func getEnvArrayStr(key string) []string {
	value := getEnv(key)

	return strings.Split(value, " | ")
}
