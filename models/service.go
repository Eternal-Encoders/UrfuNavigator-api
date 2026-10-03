package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Service struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Shape        AnyShape       `bson:"shape" json:"shape"`
	ColorSchema  *bson.ObjectID `bson:"colorSchema" json:"colorSchema"`
	IsBorder     bool           `bson:"isBorder" json:"isBorder"`
	IsFill       bool           `bson:"isFill" json:"isFill"`
}
