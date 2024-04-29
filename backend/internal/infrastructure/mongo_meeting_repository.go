package infra

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"os"
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

type MongoMeetingRepository struct {
	client           *mongo.Client
	room_collection  mongo.Collection
	event_collection mongo.Collection
}

func NewMongoMeetingRepository(client *mongo.Client) *MongoMeetingRepository {
	return &MongoMeetingRepository{
		client:           client,
		room_collection:  *client.Database("test-mongo").Collection("rooms"),
		event_collection: *client.Database("test-mongo").Collection("events"),
	}
}

func (m *MongoMeetingRepository) UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*meeting.Room, error) {
	collection := m.room_collection

	if upsertRoomInput.ID == nil { // Insert new room
		currentTime := time.Now().Unix()
		newRoom := meeting.Room{
			RoomID:    upsertRoomInput.RoomID,
			Capacity:  upsertRoomInput.Capacity,
			Equipment: upsertRoomInput.Equipment,
			Rules:     upsertRoomInput.Rules,
			IsDelete:  false,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		result, err := collection.InsertOne(ctx, newRoom)
		if err != nil {
			log.Printf("Failed to insert new room: %v", err)
			return nil, err
		}
		newRoom.ID = result.InsertedID.(primitive.ObjectID)

		return &newRoom, nil
	} else { // Update existing room
		id, err := primitive.ObjectIDFromHex(*upsertRoomInput.ID)
		if err != nil {
			log.Printf("Invalid ID format: %v", err)
			return nil, err
		}
		filter := bson.M{
			"_id":      id,
			"isDelete": false,
		}
		update := bson.M{
			"$set": bson.M{
				"roomID":    upsertRoomInput.RoomID,
				"capacity":  upsertRoomInput.Capacity,
				"equipment": upsertRoomInput.Equipment,
				"rules":     upsertRoomInput.Rules,
				"updatedAt": time.Now().Unix(),
			},
		}

		var updatedRoom meeting.Room
		err = collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedRoom)
		if err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection.
			if errors.Is(err, mongo.ErrNoDocuments) {
				log.Printf("Document with the given ID not found or deleted: %v", err)
				return nil, err
			} else {
				log.Printf("Failed to update new room: %v", err)
				return nil, err
			}
		}

		return &updatedRoom, nil
	}
}

func (m *MongoMeetingRepository) DeleteRoom(ctx context.Context, id string) (*meeting.Room, error) {
	collection := m.room_collection

	deleteRoomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}

	filter := bson.M{
		"_id":      deleteRoomID,
		"isDelete": false,
	}
	update := bson.M{"$set": bson.M{
		"isDelete":  true,
		"UpdatedAt": time.Now().Unix(),
	}}

	var updatedRoom meeting.Room
	err = collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedRoom)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Document with the given ID not found or has been deleted: %v", err)
			return nil, err
		} else {
			log.Printf("Failed to soft delete room: %v", err)
			return nil, err
		}
	}

	// Should the ID be released once a room was soft deleted?

	return &updatedRoom, nil

	/*
		result, err := collection.UpdateOne(ctx, filter, update)

		if result.ModifiedCount == 0 { // Check if the document was found and soft deleted
			fmt.Println("Document with the given ID not found")
			return nil, mongo.ErrNoDocuments
		} else if err != nil { // other failures
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
	*/
}

func (m *MongoMeetingRepository) QueryPaginatedRoom(ctx context.Context, first int, last int, before string, after string) (*model.RoomConnection, error) {
	collection := m.room_collection
	filter := bson.M{}

	if before != "" {
		decodedCursor, err := decodeCursor(before)
		if err != nil {
			return nil, err
		}
		filter["roomId"] = bson.M{"$lt": decodedCursor}
	}

	if after != "" {
		decodedCursor, err := decodeCursor(after)
		if err != nil {
			return nil, err
		}
		filter["roomId"] = bson.M{"$gt": decodedCursor}
	}

	options := options.Find()
	if last > 0 {
		options.SetSort(bson.D{{Key: "roomId", Value: -1}})
	}

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []*model.Room
	var actualCount int
	for cursor.Next(ctx) {
		var room model.Room
		if err := cursor.Decode(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
		actualCount++

		if first > 0 && actualCount >= first {
			if cursor.Next(ctx) {
				actualCount++
			}
			break
		}
		if last > 0 && actualCount >= last {
			if cursor.Next(ctx) {
				actualCount++
			}
			break
		}
	}

	var hasNextPage, hasPreviousPage bool
	if first != 0 && after != "" {
		hasPreviousPage = true
		hasNextPage = actualCount > first
	} else if first != 0 && after == "" {
		hasNextPage = actualCount > first
	} else if last != 0 && before != "" {
		hasPreviousPage = actualCount > last
		hasNextPage = true
	} else if last != 0 && before == "" {
		hasPreviousPage = actualCount > last
	}

	pageInfo := &model.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	edges := make([]*model.RoomEdge, len(rooms))
	for i, room := range rooms {
		edges[i] = &model.RoomEdge{
			Node:   room,
			Cursor: encodeCursor(room.RoomID),
		}
	}

	if len(rooms) > 0 {
		startCursor := encodeCursor(rooms[0].RoomID)
		pageInfo.StartCursor = &startCursor
		endCursor := encodeCursor(rooms[len(rooms)-1].RoomID)
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

func (m *MongoMeetingRepository) UpsertEvent(ctx context.Context, upsertEventInput model.UpsertEventInput) (*meeting.Event, error) {
	collection := m.event_collection

	if upsertEventInput.ID == nil { // Insert new event
		currentTime := time.Now().Unix()
		newEvent := meeting.Event{
			Title:           upsertEventInput.Title,
			Description:     upsertEventInput.Description,
			StartAt:         upsertEventInput.StartAt,
			EndAt:           upsertEventInput.EndAt,
			RoomID:          upsertEventInput.RoomID,
			ParticipantsIDs: upsertEventInput.ParticipantsIDs,
			Notes:           upsertEventInput.Notes,
			RemindAt:        upsertEventInput.RemindAt,
			IsDelete:        false,
			CreatedAt:       currentTime,
			UpdatedAt:       currentTime,
		}
		result, err := collection.InsertOne(ctx, newEvent)
		if err != nil {
			log.Printf("Failed to insert new event: %v", err)
			return nil, err
		}
		newEvent.ID = result.InsertedID.(primitive.ObjectID)

		return &newEvent, nil
	} else { // Update existing event
		id, err := primitive.ObjectIDFromHex(*upsertEventInput.ID)
		if err != nil {
			log.Printf("Invalid ID format: %v", err)
			return nil, err
		}
		filter := bson.M{
			"_id":      id,
			"isDelete": false,
		}
		update := bson.M{
			"$set": bson.M{
				"title":           upsertEventInput.Title,
				"description":     upsertEventInput.Description,
				"startAt":         upsertEventInput.StartAt,
				"endAt":           upsertEventInput.EndAt,
				"roomId":          upsertEventInput.RoomID,
				"participantsIDs": upsertEventInput.ParticipantsIDs,
				"notes":           upsertEventInput.Notes,
				"remindAt":        upsertEventInput.RemindAt,
				"updatedAt":       time.Now().Unix(),
			},
		}

		var updatedEvent meeting.Event
		err = collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedEvent)
		if err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection.
			if errors.Is(err, mongo.ErrNoDocuments) {
				log.Printf("Document with the given ID not found or deleted: %v", err)
				return nil, err
			} else {
				log.Printf("Failed to update new event: %v", err)
				return nil, err
			}
		}

		return &updatedEvent, nil
	}
}

func (m *MongoMeetingRepository) DeleteEvent(ctx context.Context, id string) (*meeting.Event, error) {
	collection := m.event_collection

	deleteEventID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}

	filter := bson.M{
		"_id":      deleteEventID,
		"isDelete": false,
	}
	update := bson.M{"$set": bson.M{
		"isDelete":  true,
		"UpdatedAt": time.Now().Unix(),
	}}

	var updatedEvent meeting.Event
	err = collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedEvent)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Document with the given ID not found or has been deleted: %v", err)
			return nil, err
		} else {
			log.Printf("Failed to soft delete event: %v", err)
			return nil, err
		}
	}

	// Should the ID be released once a room was soft deleted?
	return &updatedEvent, nil
}
