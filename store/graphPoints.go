package store

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) GetGraphPoint(id bson.ObjectID) (*models.GraphPoint, error) {
	return s.GraphPoints.GetOne(id)
}

func (s *MongoDB) GetGraphPoints(ids []bson.ObjectID) (*map[bson.ObjectID]*models.GraphPoint, error) {
	result, err := s.GraphPoints.GetMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MongoDB) FindGraphPoints(filter bson.M, limit int) ([]models.GraphPoint, error) {
	return s.GraphPoints.Find(filter, limit)
}

func (s *MongoDB) GetGraphByFloorID(floorID bson.ObjectID) (map[bson.ObjectID]models.GraphPoint, error) {
	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return nil, err
	}

	points, err := s.GraphPoints.GetMany(floor.Graph)
	if err != nil {
		return nil, err
	}

	graph := make(map[bson.ObjectID]models.GraphPoint, len(points))
	for id, point := range points {
		if point != nil {
			graph[id] = *point
		}
	}

	return graph, nil
}

func (s *MongoDB) GetExitPointsByBuilding(buildingID bson.ObjectID) ([]models.GraphPoint, error) {
	return s.FindGraphPoints(bson.M{
		"buildingId": buildingID,
		"types":        bson.M{"$in": []models.GraphPointType{models.ExitType}},
	}, 4)
}

func (s *MongoDB) GetFloorElevation(floorID bson.ObjectID) (float64, error) {
	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return 0, err
	}
	return floor.Elevation, nil
}

func (s *MongoDB) AddGraphPoint(point models.GraphPoint) (*bson.ObjectID, error) {
	return s.GraphPoints.InsertOne(point)
}

func (s *MongoDB) AddGraphPoints(points []models.GraphPoint) (*[]bson.ObjectID, error) {
	ids, err := s.GraphPoints.InsertMany(points)
	if err != nil {
		return nil, err
	}

	result := make([]bson.ObjectID, len(ids))
	for i, id := range ids {
		result[i] = *id
	}

	return &result, nil
}

func (s *MongoDB) RemoveGraphPoint(id bson.ObjectID) (*models.GraphPoint, error) {
	return s.GraphPoints.DeleteOne(id)
}

func (s *MongoDB) RemoveGraphPoints(ids []bson.ObjectID) (*map[bson.ObjectID]*models.GraphPoint, error) {
	result, err := s.GraphPoints.DeleteMany(ids)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
