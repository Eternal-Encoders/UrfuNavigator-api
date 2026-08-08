package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Room struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Shape        AnyShape        `bson:"shape"`
	PointId      *bson.ObjectID  `bson:"pointId"`
	Type         *GraphPointType `bson:"type"`
	Children     []AnyShape      `bson:"children"`
	ColorSchema  *bson.ObjectID  `bson:"colorSchema"`
	IsBorder     bool            `bson:"isBorder"`
	IsFill       bool            `bson:"isFill"`
}
