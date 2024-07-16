package main

import (
	"log"
	"os"
	"urfunavigator/index/api"
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

	if !exist {
		log.Fatal("No connection uri")
	}
	if !portExist {
		log.Fatal("No port specified")
	}

	store := store.Connect(uri)
	api := api.NewAPI(port, store)

	defer store.Disconnect()

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
