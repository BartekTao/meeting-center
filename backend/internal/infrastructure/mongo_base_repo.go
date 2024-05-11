package infra

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository[T any] struct{}

func (r *BaseRepository[T]) getByID(ctx context.Context, collection *mongo.Collection, id string) (*T, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID:", err)
		return nil, err
	}
	filter := bson.M{"_id": _id}
	var result T
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("Failed to decode document: %v", err)
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) softDelete(ctx context.Context, collection *mongo.Collection, id string) (*T, error) {
	deleteID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID:", err)
		return nil, err
	}

	filter := bson.M{
		"_id":      deleteID,
		"isDelete": false,
	}
	update := bson.M{
		"$set": bson.M{
			"isDelete":  true,
			"UpdatedAt": time.Now().Unix(),
		},
	}
	var result T
	err = collection.FindOneAndUpdate(
		ctx, filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Document with the given ID not found or already marked as deleted: %v", err)
			return nil, err
		} else {
			log.Printf("Failed to soft delete document: %v", err)
			return nil, err
		}
	}

	return &result, nil
}
