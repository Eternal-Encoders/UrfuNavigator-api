package models

type InstituteGps struct {
	Centre float64 `json:"centre"`
	Floor  int64   `json:"floor"`
}

type LinearGps struct {
	B1 float64 `json:"b1"`
	B2 float64 `json:"b2"`
	A  float64 `json:"a"`
}

type PointGps[T any] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

type StabilForceGps struct {
	Point PointGps[float64] `json:"point"`
	Force PointGps[float64] `json:"force"`
}

type FloorGps struct {
	Linear PointGps[LinearGps] `json:"linear"`
	Forces []StabilForceGps    `json:"forces"`
}
