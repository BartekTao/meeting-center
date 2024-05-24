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
		if err == mongo.ErrNoDocuments {
			log.Println("No document was found")
			return nil, nil
		} else {
			log.Printf("Failed to decode document: %v", err)
			return nil, err
		}
	}
	return &result, nil
}

func (r *BaseRepository[T]) findOneByFilter(
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
) (*T, error) {
	var result T
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Failed to decode document: %v", err)
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) findAllByFilter(
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
) ([]*T, error) {
	var results []*T
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to find documents: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var element T
		if err := cursor.Decode(&element); err != nil {
			log.Printf("Failed to decode document: %v", err)
			return nil, err
		}
		results = append(results, &element)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return results, nil
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

func (r *BaseRepository[T]) queryPaginated(
	ctx context.Context,
	collection *mongo.Collection,
	skip int, limit int,
	filter bson.M,
	sort bson.D,
) ([]*T, error) {
	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*T
	for cur.Next(ctx) {
		var result T
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *BaseRepository[T]) updateOne(
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
	update bson.M,
) (*mongo.UpdateResult, error) {
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Document with the given ID not found: %v", err)
			return nil, err
		} else {
			log.Printf("Failed to update: %v", err)
			return nil, err
		}
	}
	return res, err
}
