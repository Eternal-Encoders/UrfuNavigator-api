package models

import "go.mongodb.org/mongo-driver/v2/bson"

type BuildingGps struct {
	CentreAltitude float64       `bson:"centreAltitude"`
	FloorId        bson.ObjectID `bson:"floorId"`
}

type Building struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Floors       []bson.ObjectID `bson:"floors"`
	Url          string          `bson:"url"`
	Latitude     float64         `bson:"latitude"`
	Longitude    float64         `bson:"longitude"`
	Icon         string          `bson:"icon"`
	ColorSchemes []bson.ObjectID `bson:"colorSchemes"`
	Gps          []BuildingGps   `bson:"gps"`
}
