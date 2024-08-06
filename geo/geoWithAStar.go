package geo

import (
	"errors"
	"urfunavigator/index/models"
)

type GeoWithAStar struct{}

func ConnectGeoService() *GeoWithAStar {
	return &GeoWithAStar{}
}

func (g *GeoWithAStar) FindPath(
	start models.GraphPoint,
	end models.GraphPoint,
	floor_fn getFloor,
	stairs_fn getStairs,
	enters_fn getEnters,
) (map[string]map[int][][]models.GraphPoint, error) {
	if start.Institute == end.Institute {
		res, err := AStarSearch(start, end, floor_fn, stairs_fn)
		return map[string]map[int][][]models.GraphPoint{
			start.Institute: res,
		}, err
	}

	startEnters, startErr := enters_fn(start.Institute)
	endEnters, endErr := enters_fn(end.Institute)
	if startErr != nil {
		return nil, startErr
	}
	if len(startEnters) <= 0 {
		return nil, errors.New("can't find exit in start institute")
	}
	if endErr != nil {
		return nil, endErr
	}
	if len(endEnters) <= 0 {
		return nil, errors.New("can't find exit in end institute")
	}

	startPath, startPathErr := AStarSearch(start, endEnters[0], floor_fn, stairs_fn)
	endPath, endPathErr := AStarSearch(startEnters[0], end, floor_fn, stairs_fn)

	return map[string]map[int][][]models.GraphPoint{
		start.Institute: startPath,
		end.Institute:   endPath,
	}, errors.Join(startPathErr, endPathErr)
}
