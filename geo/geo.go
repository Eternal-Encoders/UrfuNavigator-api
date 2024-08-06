package geo

import "urfunavigator/index/models"

type getFloor func(floor int, institue string) (map[string]models.GraphPoint, error)
type getStairs func(institue string) (map[string]models.StairParsed, error)
type getEnters func(institute string) ([]models.GraphPoint, error)

type Geo interface {
	FindPath(
		start models.GraphPoint,
		end models.GraphPoint,
		floor_fn getFloor,
		stairs_fn getStairs,
		enters_fn getEnters,
	) (map[string]map[int][][]models.GraphPoint, error)
}
