package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BaseDBSchema struct {
	Id              bson.ObjectID `bson:"_id" json:"id"`
	DisplayableName string        `bson:"displayableName" json:"displayableName"`
	CreatedAt       time.Time     `bson:"createdAt" json:"createdAt"`
	UpdateAt        time.Time     `bson:"updateAt" json:"updatedAt"`
	Author          bson.ObjectID `bson:"author" json:"author"`
	LastUpdatedBy   bson.ObjectID `bson:"lastUpdatedBy" json:"lastUpdatedBy"`
}

func (b BaseDBSchema) GetID() bson.ObjectID {
	return b.Id
}

type GraphPointType string

const (
	CorridorsType  GraphPointType = "corridor"
	AuditoriumType GraphPointType = "auditorium"
	DinningType    GraphPointType = "dinning"
	ExitType       GraphPointType = "exit"
	StairType      GraphPointType = "stair"
	ToiletMType    GraphPointType = "toilet-m"
	ToiletWType    GraphPointType = "toilet-w"
	CafeType       GraphPointType = "cafe"
	VendingType    GraphPointType = "vending"
	CoworkingType  GraphPointType = "coworking"
	AtmType        GraphPointType = "atm"
	WardrobeType   GraphPointType = "wardrobe"
	PrintType      GraphPointType = "print"
	DeaneryType    GraphPointType = "deanery"
	StudentsType   GraphPointType = "students"
	OtherType      GraphPointType = "other"
)
