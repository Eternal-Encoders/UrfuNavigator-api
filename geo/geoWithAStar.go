package geo

import (
	"errors"
	"urfunavigator/index/logger"
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GeoWithAStar struct{}

func ConnectGeoService() *GeoWithAStar {
	return &GeoWithAStar{}
}

func (g *GeoWithAStar) FindPath(
	start models.GraphPoint,
	end models.GraphPoint,
	graphFn models.GetFloorGraph,
	elevationFn models.GetFloorElevation,
	exitsFn models.GetBuildingExits,
	getPoint func(bson.ObjectID) (*models.GraphPoint, error),
) (map[bson.ObjectID]map[bson.ObjectID][][]models.GraphPoint, error) {
	env := newPathEnvironment(graphFn, elevationFn, getPoint)

	if start.BuildingId == end.BuildingId {
		logger.Debug("path search within single building",
			"building_id", start.BuildingId.Hex(),
			"from", start.GetID().Hex(),
			"to", end.GetID().Hex(),
		)

		path, err := astarSearch(env, start, end)
		if err != nil {
			logger.Error("single-building path search failed",
				"building_id", start.BuildingId.Hex(),
				"from", start.GetID().Hex(),
				"to", end.GetID().Hex(),
				"err", err,
			)
			return nil, err
		}

		return map[bson.ObjectID]map[bson.ObjectID][][]models.GraphPoint{
			start.BuildingId: segmentPathByFloor(path),
		}, nil
	}

	logger.Info("path search across buildings",
		"start_building", start.BuildingId.Hex(),
		"end_building", end.BuildingId.Hex(),
		"from", start.GetID().Hex(),
		"to", end.GetID().Hex(),
	)

	startExits, err := exitsFn(start.BuildingId)
	if err != nil {
		return nil, err
	}
	if len(startExits) == 0 {
		return nil, errors.New("can't find exit in start building")
	}

	endExits, err := exitsFn(end.BuildingId)
	if err != nil {
		return nil, err
	}
	if len(endExits) == 0 {
		return nil, errors.New("can't find exit in end building")
	}

	if err := env.seed(append(startExits, endExits...)...); err != nil {
		return nil, err
	}

	startExit, err := selectNearestExit(start, startExits, env)
	if err != nil {
		return nil, err
	}

	endExit, err := selectNearestExit(end, endExits, env)
	if err != nil {
		return nil, err
	}

	startPath, startErr := astarSearch(env, start, startExit)
	endPath, endErr := astarSearch(env, endExit, end)
	if err := errors.Join(startErr, endErr); err != nil {
		logger.Error("cross-building path search failed",
			"start_building", start.BuildingId.Hex(),
			"end_building", end.BuildingId.Hex(),
			"start_exit", startExit.GetID().Hex(),
			"end_exit", endExit.GetID().Hex(),
			"err", err,
		)
		return nil, err
	}

	logger.Debug("cross-building path search completed",
		"start_building", start.BuildingId.Hex(),
		"end_building", end.BuildingId.Hex(),
		"start_exit", startExit.GetID().Hex(),
		"end_exit", endExit.GetID().Hex(),
	)

	return map[bson.ObjectID]map[bson.ObjectID][][]models.GraphPoint{
		start.BuildingId: segmentPathByFloor(startPath),
		end.BuildingId:   segmentPathByFloor(endPath),
	}, nil
}

var _ models.Geo = (*GeoWithAStar)(nil)
