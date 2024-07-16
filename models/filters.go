package models

type PointsFilters struct {
	Field   string
	Value   any
	IsExist bool
}

func GetPointsFilter(field string, value any, isExist bool) PointsFilters {
	return PointsFilters{
		Field:   field,
		Value:   value,
		IsExist: isExist,
	}
}
