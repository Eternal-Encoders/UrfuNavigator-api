package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetBuilding(id bson.ObjectID) (*models.Building, error) {
	return s.Buildings.GetOne(id)
}

func (s *MongoDB) GetBuildings(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Building, error) {
	result, err := s.Buildings.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) ListBuildings() ([]models.Building, error) {
	return s.Buildings.Find(bson.M{}, 0)
}

func (s *MongoDB) AddBuilding(building models.Building) (*bson.ObjectID, error) {
	return s.Buildings.InsertOne(building)
}

func (s *MongoDB) RemoveBuilding(id bson.ObjectID) (*models.Building, error) {
	building, err := s.Buildings.GetOne(id)
	if err != nil {
		return nil, err
	}

	if err := s.deleteBuildingContents(building); err != nil {
		return nil, err
	}

	return s.Buildings.DeleteOne(id)
}
