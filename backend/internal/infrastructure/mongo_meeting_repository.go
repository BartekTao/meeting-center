package infra

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
	"github.com/BartekTao/nycu-meeting-room-api/internal/meeting"
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

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Panic("Failed to ping MongoDB:", err)
	}

	log.Println("Successfully connected and pinged MongoDB.")
	return mongoClient
}

func NewMongoDBClient(ctx context.Context, cfg MongoDBConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	return client, nil
}

type MongoMeetingRepository struct {
	client *mongo.Client
}

func NewMongoMeetingRepository(client *mongo.Client) *MongoMeetingRepository {
	return &MongoMeetingRepository{
		client: client,
	}
}

func (m *MongoMeetingRepository) UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*meeting.Room, error) {
	panic(fmt.Errorf("not implemented CreateRoom - mongo repo"))
}
