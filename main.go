// @title UrfuNavigator API
// @version 1.0
// @description Campus navigation API: floors, buildings, points, search, routing, and admin CRUD.
// @host 127.0.0.1:5000
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token. Example: Bearer eyJhbGciOiJIUzI1NiIs...
//
//go:generate go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d . -o ./docs --parseInternal
package main

import (
	"os"
	"urfunavigator/index/api"
	"urfunavigator/index/auth"
	"urfunavigator/index/config"
	"urfunavigator/index/geo"
	"urfunavigator/index/logger"
	"urfunavigator/index/models"
	"urfunavigator/index/object"
	"urfunavigator/index/store"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	logger.Init(defaultEnv("LOG_LEVEL", "info"), defaultEnv("LOG_FORMAT", "text"))
	if _, err := os.Stat(".env"); err != nil {
		logger.Warn("no .env file loaded, using system environment only")
	}

	cfg := config.New()
	logger.Info("configuration loaded", "mode", cfg.Mode, "port", cfg.Port, "log_level", cfg.LogLevel)

	storeConn := store.Connect(cfg.DbUri, cfg.DbCollection)
	objectStore := object.Connect(
		cfg.S3Endpoint,
		cfg.S3Access,
		cfg.S3Secret,
		cfg.S3bucketName,
		object.ConnectOptions{
			RequestTimeout: cfg.S3RequestTimeout,
			PresignTTL:     cfg.S3PresignTTL,
			MaxRetries:     cfg.S3MaxRetries,
		},
	)
	geoService := geo.ConnectGeoService()

	server := api.NewAPI(cfg.DefaultPath, cfg.Port, cfg.Cors)

	defer func() {
		if err := storeConn.Disconnect(); err != nil {
			logger.Error("failed to disconnect mongodb", "err", err)
		}
		objectStore.Disconnect()
	}()

	dataService := models.DataService{
		Store:       storeConn,
		ObjectStore: objectStore,
		GeoService:  geoService,
	}

	logger.Info("starting http server", "port", cfg.Port, "base_path", cfg.DefaultPath)
	if err := server.Run(dataService, auth.Config{
		Secret:    []byte(cfg.JWTSecret),
		ExpiresIn: cfg.JWTExpiresIn,
	}); err != nil {
		logger.Fatal("http server stopped with error", "err", err)
	}
}

func defaultEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
