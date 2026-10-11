package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetBuildingColorSchema(id bson.ObjectID) (*models.BuildingColorSchema, error) {
	return s.BuildingColorSchemes.GetOne(id)
}

func (s *MongoDB) GetBuildingColorSchemes(ids []bson.ObjectID) (*map[bson.ObjectID]*models.BuildingColorSchema, error) {
	result, err := s.BuildingColorSchemes.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) AddBuildingColorSchema(schema models.BuildingColorSchema) (*bson.ObjectID, error) {
	return s.BuildingColorSchemes.InsertOne(schema)
}

func (s *MongoDB) AddBuildingColorSchemes(schemas []models.BuildingColorSchema) (*[]*bson.ObjectID, error) {
	ids, err := s.BuildingColorSchemes.InsertMany(schemas)
	if err != nil {
		return nil, err
	}
	return &ids, nil
}

func (s *MongoDB) RemoveBuildingColorSchema(id bson.ObjectID) (*models.BuildingColorSchema, error) {
	return s.BuildingColorSchemes.DeleteOne(id)
}

func (s *MongoDB) RemoveBuildingColorSchemes(ids []bson.ObjectID) (*map[bson.ObjectID]*models.BuildingColorSchema, error) {
	result, err := s.BuildingColorSchemes.DeleteMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
