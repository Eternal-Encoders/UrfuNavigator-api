package utils

import (
	"math"
	"urfunavigator/index/models"
)

func GraphCost(a models.GraphPoint, b models.GraphPoint) float64 {
	if a.StairId != nil || b.StairId != nil {
		return 1000
	}
	return 10
}

func Heuristic(a models.GraphPoint, b models.GraphPoint) float64 {
	if a.Floor == b.Floor {
		return math.Pow((math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2)), 0.5)
	}
	return -(2000/math.Abs(float64(a.Floor-b.Floor)) + 1)
}

func GetNeighbours(current models.GraphPoint, stairs map[string]models.StairParsed) []models.QueueItem {
	result := []models.QueueItem{}

	for _, link := range current.Links {
		result = append(result, models.QueueItem{
			Value: models.QueueItemValue{
				Id:    link,
				Floor: current.Floor,
			},
			Priority: 0,
		})
	}

	if stair, stairExist := stairs[current.Id]; stairExist {
		for _, link := range stair.Links {
			result = append(result, models.QueueItem{
				Value: models.QueueItemValue{
					Id:    link.Id,
					Floor: link.Floor,
				},
				Priority: 0,
			})
		}
	}

	return result
}

func RestorePath(path map[string]models.GraphPoint, start models.GraphPoint, end models.GraphPoint) map[int][][]models.GraphPoint {
	result := map[int][][]models.GraphPoint{
		end.Floor: {{}},
	}

	current := end
	prevFloor := -1000

	for current.Id != start.Id {
		lastPathIndex := 0
		if len(result[current.Floor]) == 0 {
			result[current.Floor] = [][]models.GraphPoint{{}}
		} else {
			lastPathIndex = len(result[current.Floor]) - 1
		}

		result[current.Floor][lastPathIndex] = append(result[current.Floor][lastPathIndex], current)

		prevFloor = current.Floor
		current = path[current.Id]

		if prevFloor != current.Floor {
			_, isFloorExist := result[current.Floor]
			if isFloorExist {
				result[current.Floor] = append(result[current.Floor], []models.GraphPoint{})
			} else {
				result[current.Floor] = [][]models.GraphPoint{}
			}
		}
	}

	result[current.Floor][len(result[current.Floor])-1] = append(result[current.Floor][len(result[current.Floor])-1], start)

	return result
}
