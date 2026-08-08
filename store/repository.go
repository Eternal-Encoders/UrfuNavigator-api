package store

import (
	"context"
	"fmt"
	"time"

	"urfunavigator/index/logger"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const queryTimeout = 10 * time.Second

type Identifiable interface {
	GetID() bson.ObjectID
}

type Repository[T Identifiable] struct {
	collection *mongo.Collection
	name       string
}

func NewRepository[T Identifiable](db *mongo.Database, collectionName string) Repository[T] {
	return Repository[T]{
		collection: db.Collection(collectionName),
		name:       collectionName,
	}
}

func (r Repository[T]) GetOne(id bson.ObjectID) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var doc T
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r Repository[T]) FindOne(filter bson.M) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var doc T
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r Repository[T]) GetMany(ids []bson.ObjectID) (map[bson.ObjectID]*T, error) {
	result := make(map[bson.ObjectID]*T, len(ids))
	for _, id := range ids {
		result[id] = nil
	}

	if len(ids) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			logger.Warn("cannot decode document", "collection", r.name, "err", err)
			continue
		}

		id := doc.GetID()
		result[id] = &doc
	}

	return result, nil
}

func (r Repository[T]) Find(filter bson.M, limit int) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []T
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r Repository[T]) InsertOne(doc T) (*bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return objectIDFromInsert(res.InsertedID)
}

func (r Repository[T]) InsertMany(docs []T) ([]*bson.ObjectID, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	docsAny := make([]any, len(docs))
	for i := range docs {
		docsAny[i] = docs[i]
	}

	res, err := r.collection.InsertMany(ctx, docsAny)
	if err != nil {
		return nil, err
	}

	ids := make([]*bson.ObjectID, len(res.InsertedIDs))
	for i, raw := range res.InsertedIDs {
		id, err := objectIDFromInsert(raw)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

func (r Repository[T]) DeleteOne(id bson.ObjectID) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var doc T
	if err := r.collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r Repository[T]) DeleteMany(ids []bson.ObjectID) (map[bson.ObjectID]*T, error) {
	result := make(map[bson.ObjectID]*T, len(ids))
	for _, id := range ids {
		result[id] = nil
	}

	if len(ids) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var found []T
	if err := cursor.All(ctx, &found); err != nil {
		return nil, err
	}

	if len(found) == 0 {
		return result, nil
	}

	deleteResult, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	if deleteResult.DeletedCount == 0 {
		return result, nil
	}

	for i := range found {
		id := found[i].GetID()
		result[id] = &found[i]
	}

	return result, nil
}

func (r Repository[T]) ReplaceOne(id bson.ObjectID, doc T) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": id}, doc)
	return err
}

func (r Repository[T]) PushToArray(id bson.ObjectID, field string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$push": bson.M{field: value},
	})
	return err
}

func (r Repository[T]) PullFromArray(id bson.ObjectID, field string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$pull": bson.M{field: value},
	})
	return err
}

func objectIDFromInsert(raw any) (*bson.ObjectID, error) {
	switch id := raw.(type) {
	case bson.ObjectID:
		return &id, nil
	case *bson.ObjectID:
		return id, nil
	default:
		return nil, fmt.Errorf("unexpected inserted id type %T", raw)
	}
}
