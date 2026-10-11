package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ParseObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf("invalid object id %q: %w", id, err)
	}
	return objectID, nil
}

func ParseObjectIDPtr(id *string) (*bson.ObjectID, error) {
	if id == nil {
		return nil, nil
	}
	objectID, err := ParseObjectID(*id)
	if err != nil {
		return nil, err
	}
	return &objectID, nil
}
