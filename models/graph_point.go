package models

import "go.mongodb.org/mongo-driver/v2/bson"

type WeekTime struct {
	Start    int32 `bson:"start"`
	End      int32 `bson:"end"`
	IsDayOff bool  `bson:"isDayOff"`
}

type Week struct {
	Monday    WeekTime `bson:"monday"`
	Tuesday   WeekTime `bson:"tuesday"`
	Wednesday WeekTime `bson:"wednesday"`
	Thursday  WeekTime `bson:"thursday"`
	Friday    WeekTime `bson:"friday"`
	Saturday  WeekTime `bson:"saturday"`
	Sunday    WeekTime `bson:"sunday"`
}

type Transaltion struct {
	Language string `bson:"language"`
	Value    string `bson:"value"`
}

type PointName struct {
	Name         string        `bson:"name"`
	Translations []Transaltion `bson:"translations"`
}

type GraphPoint struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	BuildingId   bson.ObjectID    `bson:"buildingId"`
	FloorId      bson.ObjectID    `bson:"floorId"`
	X            float64          `bson:"x"`
	Y            float64          `bson:"y"`
	Links        []bson.ObjectID  `bson:"links"`
	Types        []GraphPointType `bson:"types"`
	Names        []PointName      `bson:"names"`
	Time         Week             `bson:"time"`
	Description  string           `bson:"description"`
	Info         string           `bson:"info"`
	IsPassFree   bool             `bson:"isPassFree"`
}
