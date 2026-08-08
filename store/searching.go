package store

import (
	"context"

	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const atlasSearchIndex = "point_search"

var atlasSearchTextPaths = bson.A{
	"description",
	"info",
	"types",
	"displayableName",
}

func (s *MongoDB) GetBySearchEngine(name string, limit int) ([]models.GraphPoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := s.Database.Collection(CollGraphPoints).Aggregate(ctx, pointSearchPipeline(name, limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GraphPoint
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func pointSearchPipeline(query string, limit int) mongo.Pipeline {
	return mongo.Pipeline{
		pointSearchStage(query),
		pointSearchMatchNamedStage(),
		pointSearchLimitStage(limit),
	}
}

func pointSearchStage(query string) bson.D {
	return bson.D{{Key: "$search", Value: bson.D{
		{Key: "index", Value: atlasSearchIndex},
		{Key: "compound", Value: bson.D{
			{Key: "must", Value: bson.A{
				bson.D{{Key: "autocomplete", Value: bson.D{
					{Key: "query", Value: query},
					{Key: "path", Value: "names"},
				}}},
			}},
			{Key: "should", Value: bson.A{
				bson.D{{Key: "text", Value: bson.D{
					{Key: "query", Value: query},
					{Key: "path", Value: atlasSearchTextPaths},
				}}},
			}},
		}},
	}}}
}

func pointSearchMatchNamedStage() bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "names.0", Value: bson.D{{Key: "$exists", Value: true}}},
	}}}
}

func pointSearchLimitStage(limit int) bson.D {
	return bson.D{{Key: "$limit", Value: limit}}
}
