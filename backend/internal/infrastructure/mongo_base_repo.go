package infra

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI string
}

func SetUpMongoDB() *mongo.Client {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("You must set the MONGO_URI environment variable")
	}
	ctx := context.Background()

	mongoClient, err := NewMongoDBClient(ctx, MongoDBConfig{URI: mongoURI})
	if err != nil {
		log.Panic(err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Panic("Failed to ping MongoDB:", err)
	}

	log.Println("Successfully connected and pinged MongoDB.")
	return mongoClient
}

func ShutdownMongoDB(mongoClient *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Panic("Failed to disconnect MongoDB:", err)
	}
	log.Println("MongoDB client disconnected successfully.")
}

func NewMongoDBClient(ctx context.Context, cfg MongoDBConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	return client, nil
}

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
