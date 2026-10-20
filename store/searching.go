package store

import (
	"context"

	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const atlasSearchIndex = "point_search"
const atlasSynonymsStore = "synonym_mapping"

var atlasSearchTextPaths = bson.A{
	"description",
	"info",
	"types",
	"displayableName",
}

func (s *MongoDB) GetBySearchEngine(search models.PointSearch) ([]models.GraphPoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := s.Database.Collection(CollGraphPoints).Aggregate(ctx, pointSearchPipeline(search))
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

func pointSearchPipeline(search models.PointSearch) mongo.Pipeline {
	stages := mongo.Pipeline{
		pointSearchStage(search.Name),
		pointSearchMatchNamedStage(),
	}

	if search.Type != nil {
		stages = append(stages, pointSearchMatchTypeStage(*search.Type))
	}

	if search.FloorId != nil || search.BuildingId != nil {
		stages = append(stages,
			pointSearchRankStage(search.FloorId, search.BuildingId),
			pointSearchRankSortStage(),
		)
	}

	stages = append(stages,
		pointSearchSkipStage(search.Offset),
		pointSearchLimitStage(search.Length),
	)
	return stages
}

func pointSearchStage(query string) bson.D {
	return bson.D{{Key: "$search", Value: bson.D{
		{Key: "index", Value: atlasSearchIndex},
		{Key: "compound", Value: bson.D{
			{Key: "should", Value: bson.A{
				bson.D{{Key: "autocomplete", Value: bson.D{
					{Key: "query", Value: query},
					{Key: "path", Value: "names.name"},
				}}},
			}},
			{Key: "should", Value: bson.A{
				bson.D{{Key: "text", Value: bson.D{
					{Key: "query", Value: query},
					{Key: "path", Value: atlasSearchTextPaths},
					{Key: "synonyms", Value: atlasSynonymsStore},
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

func pointSearchMatchTypeStage(pointType string) bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "types", Value: pointType},
	}}}
}

func pointSearchRankStage(floorID, buildingID *bson.ObjectID) bson.D {
	branches := bson.A{}
	if floorID != nil {
		branches = append(branches, bson.D{
			{Key: "case", Value: bson.D{{Key: "$eq", Value: bson.A{"$floorId", *floorID}}}},
			{Key: "then", Value: 0},
		})
	}
	if buildingID != nil {
		branches = append(branches, bson.D{
			{Key: "case", Value: bson.D{{Key: "$eq", Value: bson.A{"$buildingId", *buildingID}}}},
			{Key: "then", Value: 1},
		})
	}

	return bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "rank", Value: bson.D{
			{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: branches},
				{Key: "default", Value: 2},
			}},
		}},
	}}}
}

func pointSearchRankSortStage() bson.D {
	return bson.D{{Key: "$sort", Value: bson.D{
		{Key: "rank", Value: 1},
		{Key: "score", Value: bson.D{{Key: "$meta", Value: "searchScore"}}},
		{Key: "_id", Value: 1},
	}}}
}

func pointSearchSkipStage(offset int) bson.D {
	return bson.D{{Key: "$skip", Value: offset}}
}

func pointSearchLimitStage(limit int) bson.D {
	return bson.D{{Key: "$limit", Value: limit}}
}
