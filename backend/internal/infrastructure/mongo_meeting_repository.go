package infra

import (
	"context"
	"log"

	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
	"github.com/BartekTao/nycu-meeting-room-api/internal/meeting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI string
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

	client := m.client
	collection := client.Database("testDB").Collection("rooms")

	objectID, err := primitive.ObjectIDFromHex(*upsertRoomInput.ID)
	if err != nil {
		log.Fatal(err)
	}

	new_room := &meeting.Room{
		ID:        objectID,
		RoomID:    upsertRoomInput.RoomID,
		Capacity:  upsertRoomInput.Capacity,
		Equipment: upsertRoomInput.Equipment,
		Rules:     upsertRoomInput.Rules,
		IsDelete:  false,
		CreatedAt: 0,
		UpdatedAt: 0,
		UpdaterId: "None",
	}

	new_room_bson := bson.M{
		"ID":        objectID,
		"RoomID":    upsertRoomInput.RoomID,
		"Capacity":  upsertRoomInput.Capacity,
		"Equipment": upsertRoomInput.Equipment,
		"Rules":     upsertRoomInput.Rules,
		"IsDelete":  false,
		"CreatedAt": 0,
		"UpdatedAt": 0,
		"UpdaterId": "None",
	}

	_, err = collection.InsertOne(ctx, new_room_bson)
	if err != nil {
		return nil, err
	}

	return new_room, nil
}

func (m *MongoMeetingRepository) DeleteRoom(ctx context.Context, id string) (*meeting.Room, error) {

	client := m.client
	collection := client.Database("testDB").Collection("rooms")

	filter := bson.M{"RoomId": id}

	result, err := collection.DeleteOne(ctx, filter)

	if result.DeletedCount == 0 { // Check if the document was found and deleted
		return nil, mongo.ErrNoDocuments // Document with the given roomID not found
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}
