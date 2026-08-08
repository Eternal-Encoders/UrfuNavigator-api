package store

import (
	"context"
	"urfunavigator/index/logger"
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database

	Users                Repository[models.User]
	Buildings            Repository[models.Building]
	BuildingColorSchemes Repository[models.BuildingColorSchema]
	Floors               Repository[models.Floor]
	GraphPoints          Repository[models.GraphPoint]
	Rooms                Repository[models.Room]
	Services             Repository[models.Service]
}

func Connect(uri string, collection string) *MongoDB {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		logger.Fatal("failed to connect to mongodb", "err", err)
	}

	db := client.Database(collection)
	logger.Info("connected to mongodb", "database", collection)

	return &MongoDB{
		Client:               client,
		Database:             db,
		Users:                NewRepository[models.User](db, CollUsers),
		Buildings:            NewRepository[models.Building](db, CollBuildings),
		BuildingColorSchemes: NewRepository[models.BuildingColorSchema](db, CollBuildingColorSchemes),
		Floors:               NewRepository[models.Floor](db, CollFloors),
		GraphPoints:          NewRepository[models.GraphPoint](db, CollGraphPoints),
		Rooms:                NewRepository[models.Room](db, CollRooms),
		Services:             NewRepository[models.Service](db, CollServices),
	}
}

func (s *MongoDB) Disconnect() error {
	return s.Client.Disconnect(context.TODO())
}
