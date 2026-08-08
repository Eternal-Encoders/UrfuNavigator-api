package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Store interface {
	// Users
	GetUser(id bson.ObjectID) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetUsers(ids []bson.ObjectID) (*map[bson.ObjectID]*User, error)
	ListUsers() ([]User, error)
	AddUser(user User) (*bson.ObjectID, error)
	UpdateUser(id bson.ObjectID, user User) error
	RemoveUser(id bson.ObjectID) (*User, error)

	// Building
	GetBuilding(id bson.ObjectID) (*Building, error)
	GetBuildings(ids []bson.ObjectID) (*map[bson.ObjectID]*Building, error)
	ListBuildings() ([]Building, error)
	AddBuilding(building Building) (*bson.ObjectID, error)
	UpdateBuilding(id bson.ObjectID, building Building) error
	RemoveBuilding(id bson.ObjectID) (*Building, error)
	AddBuildingFloor(buildingID bson.ObjectID, floor Floor) (*bson.ObjectID, error)
	AddBuildingColorSchemaLink(buildingID bson.ObjectID, schema BuildingColorSchema) (*bson.ObjectID, error)
	AddBuildingGps(buildingID bson.ObjectID, gps BuildingGps) error
	RemoveBuildingFloorLink(buildingID, floorID bson.ObjectID) error
	RemoveBuildingColorSchemaLink(buildingID, schemaID bson.ObjectID) error
	RemoveBuildingGps(buildingID bson.ObjectID, floorID bson.ObjectID) error

	// Building Color Schema
	GetBuildingColorSchema(id bson.ObjectID) (*BuildingColorSchema, error)
	GetBuildingColorSchemes(ids []bson.ObjectID) (*map[bson.ObjectID]*BuildingColorSchema, error)
	ListBuildingColorSchemes() ([]BuildingColorSchema, error)
	AddBuildingColorSchema(buildingColorSchema BuildingColorSchema) (*bson.ObjectID, error)
	AddBuildingColorSchemes(buildingColorSchema []BuildingColorSchema) (*[]*bson.ObjectID, error)
	UpdateBuildingColorSchema(id bson.ObjectID, schema BuildingColorSchema) error
	RemoveBuildingColorSchema(id bson.ObjectID) (*BuildingColorSchema, error)
	RemoveBuildingColorSchemes(id []bson.ObjectID) (*map[bson.ObjectID]*BuildingColorSchema, error)

	// Floors
	GetFloor(id bson.ObjectID) (*Floor, error)
	GetFloors(ids []bson.ObjectID) (*map[bson.ObjectID]*Floor, error)
	ListFloors() ([]Floor, error)
	AddFloor(floor Floor) (*bson.ObjectID, error)
	AddFloors(floor []Floor) (*[]*bson.ObjectID, error)
	UpdateFloor(id bson.ObjectID, floor Floor) error
	RemoveFloor(id bson.ObjectID) (*Floor, error)
	RemoveFloors(id []bson.ObjectID) (*map[bson.ObjectID]*Floor, error)
	AddFloorRoom(floorID bson.ObjectID, room Room) (*bson.ObjectID, error)
	AddFloorService(floorID bson.ObjectID, service Service) (*bson.ObjectID, error)
	AddFloorGraphPoint(floorID bson.ObjectID, point GraphPoint) (*bson.ObjectID, error)
	RemoveFloorRoomLink(floorID, roomID bson.ObjectID) error
	RemoveFloorServiceLink(floorID, serviceID bson.ObjectID) error
	RemoveFloorGraphPointLink(floorID, pointID bson.ObjectID) error

	// GraphPoints
	GetGraphPoint(id bson.ObjectID) (*GraphPoint, error)
	GetGraphPoints(ids []bson.ObjectID) (*map[bson.ObjectID]*GraphPoint, error)
	ListGraphPoints() ([]GraphPoint, error)
	FindGraphPoints(filter bson.M, limit int) ([]GraphPoint, error)
	GetGraphByFloorID(floorID bson.ObjectID) (map[bson.ObjectID]GraphPoint, error)
	GetExitPointsByBuilding(buildingID bson.ObjectID) ([]GraphPoint, error)
	GetFloorElevation(floorID bson.ObjectID) (float64, error)
	AddGraphPoint(graphPoint GraphPoint) (*bson.ObjectID, error)
	AddGraphPoints(graphPoints []GraphPoint) (*[]bson.ObjectID, error)
	UpdateGraphPoint(id bson.ObjectID, point GraphPoint) error
	RemoveGraphPoint(id bson.ObjectID) (*GraphPoint, error)
	RemoveGraphPoints(ids []bson.ObjectID) (*map[bson.ObjectID]*GraphPoint, error)

	// Room
	GetRoom(id bson.ObjectID) (*Room, error)
	GetRooms(ids []bson.ObjectID) (*map[bson.ObjectID]*Room, error)
	ListRooms() ([]Room, error)
	AddRoom(room Room) (*bson.ObjectID, error)
	AddRooms(rooms []Room) (*[]*bson.ObjectID, error)
	UpdateRoom(id bson.ObjectID, room Room) error
	RemoveRoom(id bson.ObjectID) (*Room, error)
	RemoveRooms(ids []bson.ObjectID) (*map[bson.ObjectID]*Room, error)

	// Service
	GetService(id bson.ObjectID) (*Service, error)
	GetServices(ids []bson.ObjectID) (*map[bson.ObjectID]*Service, error)
	ListServices() ([]Service, error)
	AddService(service Service) (*bson.ObjectID, error)
	AddServices(service []Service) (*[]*bson.ObjectID, error)
	UpdateService(id bson.ObjectID, service Service) error
	RemoveService(id bson.ObjectID) (*Service, error)
	RemoveServices(ids []bson.ObjectID) (*map[bson.ObjectID]*Service, error)

	// Searching
	GetBySearchEngine(name string, length int) ([]GraphPoint, error)
}

type BuildingIconURL struct {
	URL       string
	ExpiresAt time.Time
}

type ObjectStore interface {
	GetBuildingIconURL(fileName string) (BuildingIconURL, error)
	CheckBuildingIcon(fileName string) (bool, error)
	RemoveBuildingIcon(fileName string) error
	PutBuildingIcon(fileName string, fileData []byte) error
}

type DataService struct {
	Store       Store
	ObjectStore ObjectStore
	GeoService  Geo
}
