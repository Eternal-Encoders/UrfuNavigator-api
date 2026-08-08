package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Service struct {
	BaseDBSchema `bson:",inline" json:",inline"`
	Shape        AnyShape       `bson:"shape"`
	ColorSchema  *bson.ObjectID `bson:"colorSchema"`
	IsBorder     bool           `bson:"isBorder"`
	IsFill       bool           `bson:"isFill"`
}
