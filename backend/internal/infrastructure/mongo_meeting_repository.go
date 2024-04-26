package infra

import (
	"context"
	"encoding/base64"
	"fmt"
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

	if upsertRoomInput.ID == nil {
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
			log.Fatalf("Failed to insert new room: %v", err)
			return nil, err
		}
		newRoom.ID = result.InsertedID.(primitive.ObjectID)

		return &newRoom, nil
	} else {
		id, err := primitive.ObjectIDFromHex(*upsertRoomInput.ID)
		if err != nil {
			log.Fatalf("Invalid ID format: %v", err)
			return nil, err
		}
		filter := bson.M{"_id": id}
		update := bson.M{
			"$set": bson.M{
				"roomID":    upsertRoomInput.RoomID,
				"capacity":  upsertRoomInput.Capacity,
				"equipment": upsertRoomInput.Equipment,
				"rules":     upsertRoomInput.Rules,
				"updatedAt": time.Now().Unix(),
			},
		}
		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatalf("Failed to update room: %v", err)
			return nil, err
		}
		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("no room found with ID %s", *upsertRoomInput.ID)
		}

		var updatedRoom meeting.Room
		if err := collection.FindOne(ctx, filter).Decode(&updatedRoom); err != nil {
			log.Fatalf("Failed to retrieve updated room: %v", err)
			return nil, err
		}

		return &updatedRoom, nil
	}
}

func (m *MongoMeetingRepository) DeleteRoom(ctx context.Context, id string) (*meeting.Room, error) {
	collection := m.room_collection

	deleteRoomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": deleteRoomID}
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
