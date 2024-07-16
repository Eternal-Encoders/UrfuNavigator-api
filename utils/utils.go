package utils

import (
	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/bson"
)

func AppendFilter(root bson.M, filter models.PointsFilters) bson.M {
	if filter.IsExist {
		root[filter.Field] = filter.Value
	}
	return root
}
