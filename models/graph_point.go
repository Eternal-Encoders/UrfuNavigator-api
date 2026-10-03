package models

import "go.mongodb.org/mongo-driver/v2/bson"

type WeekTime struct {
	Start    int32 `bson:"start" json:"start"`
	End      int32 `bson:"end" json:"end"`
	IsDayOff bool  `bson:"isDayOff" json:"isDayOff"`
}

type Week struct {
	Monday    WeekTime `bson:"monday" json:"monday"`
	Tuesday   WeekTime `bson:"tuesday" json:"tuesday"`
	Wednesday WeekTime `bson:"wednesday" json:"wednesday"`
	Thursday  WeekTime `bson:"thursday" json:"thursday"`
	Friday    WeekTime `bson:"friday" json:"friday"`
	Saturday  WeekTime `bson:"saturday" json:"saturday"`
	Sunday    WeekTime `bson:"sunday" json:"sunday"`
}

type Transaltion struct {
	Language string `bson:"language" json:"language"`
	Value    string `bson:"value" json:"value"`
}

type PointName struct {
	Name         string        `bson:"name" json:"name"`
	Translations []Transaltion `bson:"translations" json:"translations"`
}

type GraphPoint struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	BuildingId   bson.ObjectID    `bson:"buildingId" json:"buildingId"`
	FloorId      bson.ObjectID    `bson:"floorId" json:"floorId"`
	X            float64          `bson:"x" json:"x"`
	Y            float64          `bson:"y" json:"y"`
	Links        []bson.ObjectID  `bson:"links" json:"links"`
	Types        []GraphPointType `bson:"types" json:"types"`
	Names        []PointName      `bson:"names" json:"names"`
	Time         Week             `bson:"time" json:"time"`
	Description  string           `bson:"description" json:"description"`
	Info         string           `bson:"info" json:"info"`
	IsPassFree   bool             `bson:"isPassFree" json:"isPassFree"`
}
