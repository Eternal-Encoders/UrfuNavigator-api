package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetRoom(id bson.ObjectID) (*models.Room, error) {
	return s.Rooms.GetOne(id)
}

func (s *MongoDB) GetRooms(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Room, error) {
	result, err := s.Rooms.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) AddRoom(room models.Room) (*bson.ObjectID, error) {
	return s.Rooms.InsertOne(room)
}

func (s *MongoDB) AddRooms(rooms []models.Room) (*[]*bson.ObjectID, error) {
	ids, err := s.Rooms.InsertMany(rooms)
	if err != nil {
		return nil, err
	}
	return &ids, nil
}

func (s *MongoDB) RemoveRoom(id bson.ObjectID) (*models.Room, error) {
	return s.Rooms.DeleteOne(id)
}

func (s *MongoDB) RemoveRooms(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Room, error) {
	result, err := s.Rooms.DeleteMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
