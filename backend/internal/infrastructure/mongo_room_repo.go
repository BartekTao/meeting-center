package infra

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"roomID"`
	Capacity  int                `bson:"capacity"`
	Equipment []string           `bson:"equipment"`
	Rules     []string           `bson:"rules"`
	IsDelete  bool               `bson:"isDelete"`
	CreatedAt int64              `bson:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt"`
	UpdaterId string             `bson:"updaterId"`
}

type MongoRoomRepository struct {
	client         *mongo.Client
	roomCollection mongo.Collection
}

func NewRoomRepository(client *mongo.Client) domain.RoomRepository {
	return &MongoRoomRepository{
		client:         client,
		roomCollection: *client.Database("meetingCenter").Collection("rooms"),
	}
}

func (m *MongoRoomRepository) UpsertRoom(ctx context.Context, room domain.Room) (*domain.Room, error) {
	collection := m.roomCollection

	if room.ID == nil { // Insert new room
		currentTime := time.Now().Unix()
		newRoom := Room{
			RoomID:    room.RoomID,
			Capacity:  room.Capacity,
			Equipment: room.Equipment,
			Rules:     room.Rules,
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

		newDomainRoom := domain.Room{
			ID:        ptr(newRoom.ID.Hex()),
			RoomID:    newRoom.RoomID,
			Capacity:  newRoom.Capacity,
			Equipment: newRoom.Equipment,
			Rules:     newRoom.Rules,
			IsDelete:  false,
			CreatedAt: newRoom.CreatedAt,
			UpdatedAt: newRoom.UpdatedAt,
		}

		return &newDomainRoom, nil
	} else { // Update existing room
		id, err := primitive.ObjectIDFromHex(*room.ID)
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
				"roomID":    room.RoomID,
				"capacity":  room.Capacity,
				"equipment": room.Equipment,
				"rules":     room.Rules,
				"updatedAt": time.Now().Unix(),
			},
		}

		var updatedRoom domain.Room
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

func (m *MongoRoomRepository) DeleteRoom(ctx context.Context, id string) (*domain.Room, error) {
	collection := m.roomCollection

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

	var updatedRoom domain.Room
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
}

func (m *MongoRoomRepository) QueryPaginatedRoom(ctx context.Context, skip int, limit int) ([]domain.Room, error) {
	collection := m.roomCollection
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "CreatedAt", Value: 1}})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []domain.Room
	for cur.Next(ctx) {
		var result domain.Room
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func ptr(s string) *string { return &s }
