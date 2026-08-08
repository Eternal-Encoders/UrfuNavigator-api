package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetUser(id bson.ObjectID) (*models.User, error) {
	return s.Users.GetOne(id)
}

func (s *MongoDB) GetUserByLogin(login string) (*models.User, error) {
	return s.Users.FindOne(bson.M{"login": login})
}

func (s *MongoDB) GetUsers(ids []bson.ObjectID) (*map[bson.ObjectID]*models.User, error) {
	result, err := s.Users.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) AddUser(user models.User) (*bson.ObjectID, error) {
	return s.Users.InsertOne(user)
}

func (s *MongoDB) RemoveUser(id bson.ObjectID) (*models.User, error) {
	return s.Users.DeleteOne(id)
}
