package models

import "go.mongodb.org/mongo-driver/v2/bson"

type BuildingGps struct {
	CentreAltitude float64       `bson:"centreAltitude" json:"centreAltitude"`
	FloorId        bson.ObjectID `bson:"floorId" json:"floorId"`
}

type Building struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Floors       []bson.ObjectID `bson:"floors" json:"floors"`
	Url          string          `bson:"url" json:"url"`
	Latitude     float64         `bson:"latitude" json:"latitude"`
	Longitude    float64         `bson:"longitude" json:"longitude"`
	Icon         string          `bson:"icon" json:"icon"`
	ColorSchemes []bson.ObjectID `bson:"colorSchemes" json:"colorSchemes"`
	Gps          []BuildingGps   `bson:"gps" json:"gps"`
}
