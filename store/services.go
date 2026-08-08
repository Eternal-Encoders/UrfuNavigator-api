package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetService(id bson.ObjectID) (*models.Service, error) {
	return s.Services.GetOne(id)
}

func (s *MongoDB) GetServices(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Service, error) {
	result, err := s.Services.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) AddService(service models.Service) (*bson.ObjectID, error) {
	return s.Services.InsertOne(service)
}

func (s *MongoDB) AddServices(services []models.Service) (*[]*bson.ObjectID, error) {
	ids, err := s.Services.InsertMany(services)
	if err != nil {
		return nil, err
	}
	return &ids, nil
}

func (s *MongoDB) RemoveService(id bson.ObjectID) (*models.Service, error) {
	return s.Services.DeleteOne(id)
}

func (s *MongoDB) RemoveServices(ids []bson.ObjectID) (*map[bson.ObjectID]*models.Service, error) {
	result, err := s.Services.DeleteMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
