package store

import (
	"context"
	"log"
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(uri string, collection string) *MongoDB {
	dbOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), dbOptions)

	if err != nil {
		log.Fatal(err)
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(collection),
	}
}

func (s *MongoDB) Disconnect() error {
	return s.Client.Disconnect(context.TODO())
}

func (s *MongoDB) GetFloor(floor int, institute string) (models.Floor, error) {
	coll := s.Database.Collection("floors")
	filter := bson.D{
		{
			Key:   "floor",
			Value: floor,
		},
		{
			Key:   "institute",
			Value: institute,
		},
	}
	var result models.Floor
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (s *MongoDB) GetGraphPointsById(ids []string) ([]models.GraphPoint, error) {
	coll := s.Database.Collection("graph_points")
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	curs, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer curs.Close(context.TODO())

	var result []models.GraphPoint
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) GetGraph(floor int, institute string) (map[string]models.GraphPoint, error) {
	floorData, floorErr := s.GetFloor(floor, institute)
	if floorErr != nil {
		return nil, floorErr
	}

	graphData, graphErr := s.GetGraphPointsById(floorData.Graph)
	if graphErr != nil {
		return nil, graphErr
	}

	result := make(map[string]models.GraphPoint)
	for _, point := range graphData {
		result[point.Id] = point
	}
	return result, nil
}

func (s *MongoDB) GetInstitute(url string) (models.Institute, error) {
	coll := s.Database.Collection("insitutes")
	filter := bson.M{
		"url": url,
	}
	var result models.Institute
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (s *MongoDB) GetInstitutes() ([]models.Institute, error) {
	coll := s.Database.Collection("insitutes")

	curs, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer curs.Close(context.TODO())

	var result []models.Institute
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) GetInstituteIcons(ids []string) ([]models.InstituteIcon, error) {
	coll := s.Database.Collection("media")

	objIds := []primitive.ObjectID{}
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIds = append(objIds, objId)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	}
	curs, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer curs.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) GetStairs(institute string) (map[string]models.StairParsed, error) {
	coll := s.Database.Collection("stairs")
	filter := bson.M{
		"institute": institute,
	}

	curs, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer curs.Close(context.TODO())

	result := make(map[string]models.StairParsed)
	for curs.Next(context.TODO()) {
		var stair models.Stair
		if err := curs.Decode(&stair); err != nil {
			return nil, err
		}

		links, linksErr := s.GetGraphPointsById(stair.Links)
		if linksErr != nil {
			return nil, linksErr
		}
		result[stair.Id] = models.StairParsed{
			Id:    stair.Id,
			Links: links,
		}
	}

	return result, nil
}

func (s *MongoDB) GetPoints(preFilters []models.PointsFilters, limit int) ([]models.GraphPoint, error) {
	coll := s.Database.Collection("graph_points")
	filter := bson.M{}

	for _, singleFilter := range preFilters {
		filter[singleFilter.Field] = singleFilter.Value
	}

	curs, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer curs.Close(context.TODO())

	result := []models.GraphPoint{}
	for i := 0; curs.Next(context.TODO()) && i < limit; i++ {
		var point models.GraphPoint
		if err := curs.Decode(&point); err != nil {
			return nil, err
		}
		result = append(result, point)

	}

	if err := curs.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MongoDB) GetPoint(id string) (models.GraphPoint, error) {
	coll := s.Database.Collection("graph_points")
	filter := bson.M{
		"_id": id,
	}

	var result models.GraphPoint
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (s *MongoDB) GetEnters(institute string) ([]models.GraphPoint, error) {
	pointsFilter := []models.PointsFilters{
		models.GetPointsFilter("types", bson.M{"$in": []string{"exit"}}),
		models.GetPointsFilter("institute", institute),
	}

	return s.GetPoints(pointsFilter, 4)
}

func (s *MongoDB) GetBySearchEngine(name string, length int) ([]models.GraphPoint, error) {
	coll := s.Database.Collection("graph_points")

	searchStage := bson.D{
		{Key: "$search", Value: bson.D{
			{Key: "index", Value: "point_search"},
			{Key: "compound", Value: bson.D{
				{Key: "must", Value: bson.A{
					bson.D{
						{Key: "autocomplete", Value: bson.D{
							{Key: "query", Value: name},
							{Key: "path", Value: "names"},
						}},
					},
				}},
				{Key: "should", Value: bson.A{
					bson.D{
						{Key: "text", Value: bson.D{
							{Key: "query", Value: name},
							{Key: "path", Value: bson.A{
								"floor",
								"info",
								"institute",
								"types",
								"description",
							}},
						}},
					},
				}},
			}},
		}},
	}
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "names.0", Value: bson.D{
				{Key: "$exists", Value: true},
			}},
		}},
	}
	limitStage := bson.D{{Key: "$limit", Value: length}}

	curs, cursErr := coll.Aggregate(context.TODO(), mongo.Pipeline{
		searchStage,
		matchStage,
		limitStage,
	})

	if cursErr != nil {
		return nil, cursErr
	}

	defer curs.Close(context.TODO())

	var results []models.GraphPoint
	if err := curs.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
