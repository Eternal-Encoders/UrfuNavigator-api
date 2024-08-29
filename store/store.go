package store

import (
	"urfunavigator/index/models"
)

type Store interface {
	GetFloor(floor int, institute string) (models.Floor, error)
	GetGraph(floor int, institute string) (map[string]models.GraphPoint, error)
	GetGraphPointsById(ids []string) ([]models.GraphPoint, error)
	GetInstitute(url string) (models.Institute, error)
	GetInstitutes() ([]models.Institute, error)
	GetInstituteIcons(ids []string) ([]models.InstituteIcon, error)
	GetStairs(institute string) (map[string]models.StairParsed, error)
	GetPoints(preFilter []models.PointsFilters, limit int) ([]models.GraphPoint, error)
	GetPoint(id string) (models.GraphPoint, error)
	GetEnters(institute string) ([]models.GraphPoint, error)
	GetBySearchEngine(name string, length int) ([]models.GraphPoint, error)
}
