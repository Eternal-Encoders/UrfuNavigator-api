package main

import (
	"log"
	"os"
	"urfunavigator/index/api"
	"urfunavigator/index/object"
	"urfunavigator/index/store"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file. Env load from system only")
	}
}

func main() {
	uri, exist := os.LookupEnv("DATABASE_URI")
	port, portExist := os.LookupEnv("PORT")
	s3Endpoint, s3EndpointExist := os.LookupEnv("BUCKET_ENDPOINT")
	s3Access, s3AccessExist := os.LookupEnv("BUCKET_ACCESS_KEY")
	s3Secret, s3SecretExist := os.LookupEnv("BUCKET_SECRET_KEY")

	bucketName := "navigator"

	if !exist {
		log.Fatal("No connection uri")
	}
	if !portExist {
		log.Fatal("No port specified")
	}
	if !s3EndpointExist {
		log.Fatal("No s3 endpoint specified")
	}
	if !s3AccessExist {
		log.Fatal("No s3 access key specified")
	}
	if !s3SecretExist {
		log.Fatal("No s3 secret key specified")
	}

	store := store.Connect(uri)
	objectStore := object.Connect(s3Endpoint, s3Access, s3Secret, bucketName)
	api := api.NewAPI(
		port,
		store,
		objectStore,
		"http://localhost:3000, https://how-to-navigate.ru, https://how-to-navigate.ru:2053/",
	)

	defer store.Disconnect()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
