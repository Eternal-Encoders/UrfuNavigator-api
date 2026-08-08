package models

import "go.mongodb.org/mongo-driver/v2/bson"

type GetFloorGraph func(floorID bson.ObjectID) (map[bson.ObjectID]GraphPoint, error)

type GetFloorElevation func(floorID bson.ObjectID) (float64, error)

type GetBuildingExits func(buildingID bson.ObjectID) ([]GraphPoint, error)

type Geo interface {
	FindPath(
		start GraphPoint,
		end GraphPoint,
		graphFn GetFloorGraph,
		elevationFn GetFloorElevation,
		exitsFn GetBuildingExits,
		getPoint func(bson.ObjectID) (*GraphPoint, error),
	) (map[bson.ObjectID]map[bson.ObjectID][][]GraphPoint, error)
}
