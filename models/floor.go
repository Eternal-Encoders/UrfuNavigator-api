package models

import "go.mongodb.org/mongo-driver/v2/bson"

type LinearGps struct {
	B1 float64 `bson:"b1"`
	B2 float64 `bson:"b2"`
	A  float64 `bson:"a"`
}

type PointGps[T any] struct {
	X T `bson:"x"`
	Y T `bson:"y"`
}

type StabilForceGps struct {
	Point PointGps[float64] `bson:"point"`
	Force PointGps[float64] `bson:"force"`
}

type FloorGps struct {
	Altitude float64             `bson:"altitude"`
	Linear   PointGps[LinearGps] `bson:"linear"`
	Forces   []StabilForceGps    `bson:"forces"`
}

type Floor struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	BuildingId   bson.ObjectID   `bson:"buildingId" json:"buildingId"`
	Elevation    float64         `bson:"elevation" json:"elevation"`
	Width        int32           `bson:"width" json:"width"`
	Height       int32           `bson:"height" json:"height"`
	Rooms        []bson.ObjectID `bson:"rooms" json:"rooms"`
	Services     []bson.ObjectID `bson:"services" json:"services"`
	Graph        []bson.ObjectID `bson:"graph" json:"graph"`
	Gps          *FloorGps       `bson:"gps" json:"gps,omitempty"`
}
