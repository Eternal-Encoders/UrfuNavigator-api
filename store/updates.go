package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) ListUsers() ([]models.User, error) {
	return s.Users.Find(bson.M{}, 0)
}

func (s *MongoDB) ListFloors() ([]models.Floor, error) {
	return s.Floors.Find(bson.M{}, 0)
}

func (s *MongoDB) ListRooms() ([]models.Room, error) {
	return s.Rooms.Find(bson.M{}, 0)
}

func (s *MongoDB) ListServices() ([]models.Service, error) {
	return s.Services.Find(bson.M{}, 0)
}

func (s *MongoDB) ListGraphPoints() ([]models.GraphPoint, error) {
	return s.GraphPoints.Find(bson.M{}, 0)
}

func (s *MongoDB) ListBuildingColorSchemes() ([]models.BuildingColorSchema, error) {
	return s.BuildingColorSchemes.Find(bson.M{}, 0)
}

func (s *MongoDB) UpdateUser(id bson.ObjectID, user models.User) error {
	user.Id = id
	return s.Users.ReplaceOne(id, user)
}

func (s *MongoDB) UpdateBuilding(id bson.ObjectID, building models.Building) error {
	building.Id = id
	return s.Buildings.ReplaceOne(id, building)
}

func (s *MongoDB) UpdateBuildingColorSchema(id bson.ObjectID, schema models.BuildingColorSchema) error {
	schema.Id = id
	return s.BuildingColorSchemes.ReplaceOne(id, schema)
}

func (s *MongoDB) UpdateFloor(id bson.ObjectID, floor models.Floor) error {
	floor.Id = id
	return s.Floors.ReplaceOne(id, floor)
}

func (s *MongoDB) UpdateRoom(id bson.ObjectID, room models.Room) error {
	room.Id = id
	return s.Rooms.ReplaceOne(id, room)
}

func (s *MongoDB) UpdateService(id bson.ObjectID, service models.Service) error {
	service.Id = id
	return s.Services.ReplaceOne(id, service)
}

func (s *MongoDB) UpdateGraphPoint(id bson.ObjectID, point models.GraphPoint) error {
	point.Id = id
	return s.GraphPoints.ReplaceOne(id, point)
}
