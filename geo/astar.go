package geo

import (
	"container/heap"
	"errors"

	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type searchNode struct {
	pointID bson.ObjectID
	fScore  float64
	index   int
}

type openSet []*searchNode

func (nodes openSet) Len() int { return len(nodes) }

func (nodes openSet) Less(i, j int) bool { return nodes[i].fScore < nodes[j].fScore }

func (nodes openSet) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
	nodes[i].index = i
	nodes[j].index = j
}

func (nodes *openSet) Push(value any) {
	node := value.(*searchNode)
	node.index = len(*nodes)
	*nodes = append(*nodes, node)
}

func (nodes *openSet) Pop() any {
	current := *nodes
	lastIndex := len(current) - 1
	node := current[lastIndex]
	current[lastIndex] = nil
	node.index = -1
	*nodes = current[:lastIndex]
	return node
}

func astarSearch(env *pathEnvironment, start, end models.GraphPoint) ([]models.GraphPoint, error) {
	if start.BuildingId != end.BuildingId {
		return nil, errors.New("building ids of start and end point must be equal")
	}

	if err := env.seed(start, end); err != nil {
		return nil, err
	}

	startID := start.GetID()
	endID := end.GetID()

	if startID == endID {
		return []models.GraphPoint{start}, nil
	}

	endElevation, err := env.elevation(end.FloorId)
	if err != nil {
		return nil, err
	}

	gScore := map[bson.ObjectID]float64{startID: 0}
	cameFrom := make(map[bson.ObjectID]bson.ObjectID)
	closed := make(map[bson.ObjectID]struct{})

	open := openSet{}
	startElevation, err := env.elevation(start.FloorId)
	if err != nil {
		return nil, err
	}

	heap.Push(&open, &searchNode{
		pointID: startID,
		fScore:  heuristic(start, end, startElevation, endElevation),
	})

	for open.Len() > 0 {
		currentID := heap.Pop(&open).(*searchNode).pointID
		if _, seen := closed[currentID]; seen {
			continue
		}
		closed[currentID] = struct{}{}

		if currentID == endID {
			return rebuildPath(env, cameFrom, startID, endID)
		}

		currentPoint, err := env.point(currentID)
		if err != nil {
			return nil, err
		}

		currentElevation, err := env.elevation(currentPoint.FloorId)
		if err != nil {
			return nil, err
		}

		neighbours, err := env.neighbours(currentPoint)
		if err != nil {
			return nil, err
		}

		for _, nextPoint := range neighbours {
			nextID := nextPoint.GetID()
			if _, seen := closed[nextID]; seen {
				continue
			}

			nextElevation, err := env.elevation(nextPoint.FloorId)
			if err != nil {
				return nil, err
			}

			stepCost := moveCost(currentPoint, nextPoint, currentElevation, nextElevation)
			tentativeG := gScore[currentID] + stepCost
			if existing, ok := gScore[nextID]; ok && tentativeG >= existing {
				continue
			}

			cameFrom[nextID] = currentID
			gScore[nextID] = tentativeG
			heap.Push(&open, &searchNode{
				pointID: nextID,
				fScore:  tentativeG + heuristic(nextPoint, end, nextElevation, endElevation),
			})
		}
	}

	return nil, errors.New("path not found")
}

func rebuildPath(
	env *pathEnvironment,
	cameFrom map[bson.ObjectID]bson.ObjectID,
	startID bson.ObjectID,
	endID bson.ObjectID,
) ([]models.GraphPoint, error) {
	path := make([]models.GraphPoint, 0)
	currentID := endID

	for {
		point, err := env.point(currentID)
		if err != nil {
			return nil, err
		}

		path = append(path, point)
		if currentID == startID {
			break
		}

		previousID, ok := cameFrom[currentID]
		if !ok {
			return nil, errors.New("path reconstruction failed")
		}
		currentID = previousID
	}

	reversePoints(path)
	return path, nil
}

func reversePoints(path []models.GraphPoint) {
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}
}

func segmentPathByFloor(path []models.GraphPoint) map[bson.ObjectID][][]models.GraphPoint {
	if len(path) == 0 {
		return map[bson.ObjectID][][]models.GraphPoint{}
	}

	result := make(map[bson.ObjectID][][]models.GraphPoint)
	currentFloorID := path[0].FloorId
	currentSegment := []models.GraphPoint{path[0]}
	result[currentFloorID] = [][]models.GraphPoint{currentSegment}

	for i := 1; i < len(path); i++ {
		point := path[i]
		if point.FloorId == currentFloorID {
			currentSegment = append(currentSegment, point)
			result[currentFloorID][len(result[currentFloorID])-1] = currentSegment
			continue
		}

		currentFloorID = point.FloorId
		currentSegment = []models.GraphPoint{point}
		result[currentFloorID] = append(result[currentFloorID], currentSegment)
	}

	return result
}

func selectNearestExit(from models.GraphPoint, exits []models.GraphPoint, env *pathEnvironment) (models.GraphPoint, error) {
	if len(exits) == 0 {
		return models.GraphPoint{}, errors.New("no exits available")
	}

	fromElevation, err := env.elevation(from.FloorId)
	if err != nil {
		return models.GraphPoint{}, err
	}

	bestExit := exits[0]
	bestExitElevation, err := env.elevation(bestExit.FloorId)
	if err != nil {
		return models.GraphPoint{}, err
	}
	bestScore := estimatePointDistance(from, bestExit, fromElevation, bestExitElevation)

	for i := 1; i < len(exits); i++ {
		exit := exits[i]
		exitElevation, err := env.elevation(exit.FloorId)
		if err != nil {
			return models.GraphPoint{}, err
		}

		score := estimatePointDistance(from, exit, fromElevation, exitElevation)
		if score < bestScore {
			bestScore = score
			bestExit = exit
		}
	}

	return bestExit, nil
}
