package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Room struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Shape        AnyShape        `bson:"shape" json:"shape"`
	PointId      *bson.ObjectID  `bson:"pointId" json:"pointId"`
	Type         *GraphPointType `bson:"type" json:"type"`
	Children     []AnyShape      `bson:"children" json:"children"`
	ColorSchema  *bson.ObjectID  `bson:"colorSchema" json:"colorSchema"`
	IsBorder     bool            `bson:"isBorder" json:"isBorder"`
	IsFill       bool            `bson:"isFill" json:"isFill"`
}
