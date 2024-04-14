package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

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
	client          *mongo.Client
	room_collection mongo.Collection
}

func NewMongoMeetingRepository(client *mongo.Client) *MongoMeetingRepository {
	return &MongoMeetingRepository{
		client:          client,
		room_collection: *client.Database("test-mongo").Collection("rooms"),
	}
}

func (m *MongoMeetingRepository) UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*meeting.Room, error) {

	collection := m.room_collection

	upsertRoomID, err := primitive.ObjectIDFromHex(*upsertRoomInput.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the update operation
	filter := bson.M{"ID": upsertRoomID}
	update := bson.M{"$set": bson.M{
		"RoomID":    upsertRoomInput.RoomID,
		"Capacity":  upsertRoomInput.Capacity,
		"Equipment": upsertRoomInput.Equipment,
		"Rules":     upsertRoomInput.Rules,
		"UpdatedAt": time.Now().Unix(),
	}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
		return nil, err
	}

	fmt.Print(result)
	var output_room meeting.Room

	if result.ModifiedCount != 0 {
		// Document with the same RoomID already exists, update succeeded
		fmt.Printf("Updated %d document(s) successfully.\n", result.ModifiedCount)
		// Retrieve the updated room document
		err = collection.FindOne(ctx, filter).Decode(&output_room)
		if err != nil {
			log.Fatalf("Failed to decode updated room document: %v", err)
			return nil, err
		}
	} else {
		// Document with the same RoomID doesn't exist, insert new room
		newRoom := bson.M{
			"ID":        upsertRoomID,
			"RoomID":    upsertRoomInput.RoomID,
			"Capacity":  upsertRoomInput.Capacity,
			"Equipment": upsertRoomInput.Equipment,
			"Rules":     upsertRoomInput.Rules,
			"IsDelete":  false,
			"CreatedAt": time.Now().Unix(),
			"UpdatedAt": time.Now().Unix(),
			"UpdaterId": "None",
		}
		_, err := collection.InsertOne(ctx, newRoom)
		if err != nil {
			log.Fatalf("Failed to insert document: %v", err)
			return nil, err
		}
		fmt.Println("New room inserted successfully.")
		err = collection.FindOne(ctx, filter).Decode(&output_room)
		if err != nil {
			log.Fatalf("Failed to decode inserted room document: %v", err)
			return nil, err
		}
	}

	return &output_room, nil
}

func (m *MongoMeetingRepository) DeleteRoom(ctx context.Context, id string) (*meeting.Room, error) {

	collection := m.room_collection

	deleteRoomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"ID": deleteRoomID}
	update := bson.M{"$set": bson.M{
		"IsDelete":  true,
		"UpdatedAt": time.Now().Unix(),
	}}

	result, err := collection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount == 0 { // Check if the document was found and soft deleted
		fmt.Println("Document with the given ID not found")
		return nil, mongo.ErrNoDocuments
	} else if err != nil { // other failures
		log.Fatalf("Failed to soft delete document: %v", err)
		return nil, err
	}

	fmt.Printf("Deleted %d document(s) successfully.\n", result.ModifiedCount)

	var deleted_room meeting.Room
	err = collection.FindOne(ctx, filter).Decode(&deleted_room)
	if err != nil {
		log.Fatalf("Failed to decode updated room document: %v", err)
		return nil, err
	}

	return &deleted_room, nil
}

func (m *MongoMeetingRepository) QueryPaginatedRoom(ctx context.Context, first int, after, before string) (*model.RoomConnection, error) {

	filter := bson.M{}
	if after != "" {
		decodedCursor, err := decodeCursor(after)
		if err != nil {
			return nil, err
		}
		filter["_id"] = bson.M{"$gt": decodedCursor}
	}

	cursor, err := m.client.Database("test-mongo").Collection("rooms").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []*model.Room
	for cursor.Next(ctx) {
		var room model.Room
		if err := cursor.Decode(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
		if len(rooms) >= first {
			break
		}
	}

	edges := make([]*model.RoomEdge, len(rooms))
	for i, room := range rooms {
		edges[i] = &model.RoomEdge{
			Node:   room,
			Cursor: encodeCursor(room.ID),
		}
	}
	hasNextPage := len(rooms) > first

	pageInfo := &model.PageInfo{
		HasNextPage: hasNextPage,
	}

	if len(rooms) > 0 {
		startCursor := encodeCursor(rooms[0].ID)
		pageInfo.StartCursor = &startCursor
		endCursor := encodeCursor(rooms[len(rooms)-1].ID)
		pageInfo.EndCursor = &endCursor
	}

	roomConnection := &model.RoomConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
	return roomConnection, nil
}

func encodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func decodeCursor(cursorString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(cursorString)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
