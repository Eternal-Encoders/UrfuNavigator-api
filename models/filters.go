package models

type PointsFilters struct {
	Field string
	Value any
}

func GetPointsFilter(field string, value any) PointsFilters {
	return PointsFilters{
		Field: field,
		Value: value,
	}
}
