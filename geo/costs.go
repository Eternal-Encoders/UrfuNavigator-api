package geo

import (
	"math"
	"urfunavigator/index/models"
)

const stairMoveCost = 1000.0

func isStairPoint(point models.GraphPoint) bool {
	for _, pointType := range point.Types {
		if pointType == models.StairType {
			return true
		}
	}
	return false
}

func moveCost(from, to models.GraphPoint, elevationFrom, elevationTo float64) float64 {
	distance := math.Hypot(from.X-to.X, from.Y-to.Y)
	if distance < 1 {
		distance = 1
	}

	if from.FloorId != to.FloorId || elevationFrom != elevationTo {
		return stairMoveCost + distance
	}

	if isStairPoint(from) || isStairPoint(to) {
		return stairMoveCost
	}

	return distance
}

func estimatePointDistance(from, to models.GraphPoint, elevationFrom, elevationTo float64) float64 {
	elevationDiff := math.Abs(elevationFrom - elevationTo)
	return math.Hypot(from.X-to.X, from.Y-to.Y) + elevationDiff*stairMoveCost
}

func heuristic(from, to models.GraphPoint, elevationFrom, elevationTo float64) float64 {
	return estimatePointDistance(from, to, elevationFrom, elevationTo)
}
