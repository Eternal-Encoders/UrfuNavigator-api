package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetFloor(id bson.ObjectID) (*models.Floor, error) {
	return s.Floors.GetOne(id)
}

func (s *MongoDB) GetFloors(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Floor, error) {
	result, err := s.Floors.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) AddFloor(floor models.Floor) (*bson.ObjectID, error) {
	return s.Floors.InsertOne(floor)
}

func (s *MongoDB) AddFloors(floors []models.Floor) (*[]*bson.ObjectID, error) {
	ids, err := s.Floors.InsertMany(floors)
	if err != nil {
		return nil, err
	}
	return &ids, nil
}

func (s *MongoDB) RemoveFloor(id bson.ObjectID) (*models.Floor, error) {
	floor, err := s.Floors.GetOne(id)
	if err != nil {
		return nil, err
	}

	if err := s.deleteFloorContents(floor); err != nil {
		return nil, err
	}

	return s.Floors.DeleteOne(id)
}

func (s *MongoDB) RemoveFloors(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Floor, error) {
	floors, err := s.Floors.GetMany(ids)
	if err != nil {
		return nil, err
	}

	if err := s.deleteFloorsContents(floors); err != nil {
		return nil, err
	}

	result, err := s.Floors.DeleteMany(ids)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
