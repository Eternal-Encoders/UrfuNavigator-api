package main

import (
	"log"
	"os"
	"urfunavigator/index/api"
	"urfunavigator/index/geo"
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
	collection, collectionExist := os.LookupEnv("DATABASE_COLLECTION")
	port, portExist := os.LookupEnv("PORT")
	cors, corsExist := os.LookupEnv("CORS")
	s3Endpoint, s3EndpointExist := os.LookupEnv("BUCKET_ENDPOINT")
	s3Access, s3AccessExist := os.LookupEnv("BUCKET_ACCESS_KEY")
	s3Secret, s3SecretExist := os.LookupEnv("BUCKET_SECRET_KEY")
	bucketName, bucketNameExist := os.LookupEnv("BUCKET_NAME")

	if !exist {
		file, fileExist := os.LookupEnv("DATABASE_URI_FILE")

		if !fileExist {
			log.Fatal("No connection uri")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		uri = string(data)
	}
	if !collectionExist {
		log.Fatal("No collection specified")
	}
	if !portExist {
		log.Fatal("No port specified")
	}
	if !corsExist {
		log.Fatal("Cors policy not specified")
	}
	if !s3EndpointExist {
		log.Fatal("No s3 endpoint specified")
	}
	if !s3AccessExist {
		file, fileExist := os.LookupEnv("BUCKET_ACCESS_KEY_FILE")

		if !fileExist {
			log.Fatal("No s3 access key specified")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		s3Access = string(data)
	}
	if !s3SecretExist {
		file, fileExist := os.LookupEnv("BUCKET_SECRET_KEY_FILE")

		if !fileExist {
			log.Fatal("No s3 secret key specified")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		s3Secret = string(data)
	}
	if !bucketNameExist {
		log.Fatal("No s3 bucket specified")
	}

	store := store.Connect(uri, collection)
	objectStore := object.Connect(s3Endpoint, s3Access, s3Secret, bucketName)
	geoService := geo.ConnectGeoService()
	api := api.NewAPI(
		"/api",
		port,
		store,
		objectStore,
		geoService,
		cors,
	)

	defer store.Disconnect()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
