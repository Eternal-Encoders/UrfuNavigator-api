package geo

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type pathEnvironment struct {
	graphFn      models.GetFloorGraph
	elevationFn  models.GetFloorElevation
	getPoint     func(bson.ObjectID) (*models.GraphPoint, error)

	floorGraphs   map[bson.ObjectID]map[bson.ObjectID]models.GraphPoint
	floorElevations map[bson.ObjectID]float64
	points        map[bson.ObjectID]models.GraphPoint
}

func newPathEnvironment(
	graphFn models.GetFloorGraph,
	elevationFn models.GetFloorElevation,
	getPoint func(bson.ObjectID) (*models.GraphPoint, error),
) *pathEnvironment {
	return &pathEnvironment{
		graphFn:         graphFn,
		elevationFn:     elevationFn,
		getPoint:        getPoint,
		floorGraphs:     make(map[bson.ObjectID]map[bson.ObjectID]models.GraphPoint),
		floorElevations: make(map[bson.ObjectID]float64),
		points:          make(map[bson.ObjectID]models.GraphPoint),
	}
}

func (env *pathEnvironment) registerPoint(point models.GraphPoint) {
	env.points[point.GetID()] = point
}

func (env *pathEnvironment) seed(points ...models.GraphPoint) error {
	for _, point := range points {
		env.registerPoint(point)
		if _, err := env.ensureFloorGraph(point.FloorId); err != nil {
			return err
		}
	}
	return nil
}

func (env *pathEnvironment) ensureFloorGraph(floorID bson.ObjectID) (map[bson.ObjectID]models.GraphPoint, error) {
	if graph, ok := env.floorGraphs[floorID]; ok {
		return graph, nil
	}

	graph, err := env.graphFn(floorID)
	if err != nil {
		return nil, err
	}

	for id, point := range graph {
		env.points[id] = point
	}

	env.floorGraphs[floorID] = graph
	return graph, nil
}

func (env *pathEnvironment) elevation(floorID bson.ObjectID) (float64, error) {
	if value, ok := env.floorElevations[floorID]; ok {
		return value, nil
	}

	value, err := env.elevationFn(floorID)
	if err != nil {
		return 0, err
	}

	env.floorElevations[floorID] = value
	return value, nil
}

func (env *pathEnvironment) point(pointID bson.ObjectID) (models.GraphPoint, error) {
	if point, ok := env.points[pointID]; ok {
		return point, nil
	}

	point, err := env.getPoint(pointID)
	if err != nil {
		return models.GraphPoint{}, err
	}

	env.registerPoint(*point)
	if _, err := env.ensureFloorGraph(point.FloorId); err != nil {
		return models.GraphPoint{}, err
	}

	return *point, nil
}

func (env *pathEnvironment) neighbours(point models.GraphPoint) ([]models.GraphPoint, error) {
	if _, err := env.ensureFloorGraph(point.FloorId); err != nil {
		return nil, err
	}

	neighbours := make([]models.GraphPoint, 0, len(point.Links))
	for _, linkID := range point.Links {
		linked, err := env.point(linkID)
		if err != nil {
			continue
		}

		if _, err := env.ensureFloorGraph(linked.FloorId); err != nil {
			return nil, err
		}

		neighbours = append(neighbours, linked)
	}

	return neighbours, nil
}
