package infra

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/common"
	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Room struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Capacity   int                `bson:"capacity"`
	Equipments []domain.Equipment `bson:"equipments"`
	Rules      []domain.Rule      `bson:"rules"`
	IsDelete   bool               `bson:"isDelete"`
	CreatedAt  int64              `bson:"createdAt"`
	CreatorID  primitive.ObjectID `bson:"creatorID"`
	UpdatedAt  int64              `bson:"updatedAt"`
	UpdaterID  primitive.ObjectID `bson:"updaterID"`
}

type mongoRoomRepository struct {
	BaseRepository[Room]
	client         *mongo.Client
	roomCollection *mongo.Collection
}

func NewMongoRoomRepository(client *mongo.Client) domain.RoomRepository {
	return &mongoRoomRepository{
		client:         client,
		roomCollection: client.Database("meetingCenter").Collection("rooms"),
	}
}

func (m *mongoRoomRepository) Upsert(ctx context.Context, room domain.Room) (*domain.Room, error) {
	collection := m.roomCollection
	updaterID, err := primitive.ObjectIDFromHex(room.UpdaterID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, err
	}

	if room.ID == nil { // Insert new room

		currentTime := time.Now().Unix()
		newRoom := Room{
			Name:       room.Name,
			Capacity:   room.Capacity,
			Equipments: room.Equipments,
			Rules:      room.Rules,
			IsDelete:   false,
			CreatedAt:  currentTime,
			CreatorID:  updaterID,
			UpdatedAt:  currentTime,
			UpdaterID:  updaterID,
		}
		result, err := collection.InsertOne(ctx, newRoom)
		if err != nil {
			log.Printf("Failed to insert new room: %v", err)
			return nil, err
		}
		newRoom.ID = result.InsertedID.(primitive.ObjectID)

		return ToDomainRoom(&newRoom), nil
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
				"name":       room.Name,
				"capacity":   room.Capacity,
				"equipments": room.Equipments,
				"rules":      room.Rules,
				"updatedAt":  time.Now().Unix(),
				"updaterID":  updaterID,
			},
		}

		var updatedRoom Room
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

		return ToDomainRoom(&updatedRoom), nil
	}
}

func (m *mongoRoomRepository) Delete(ctx context.Context, id string) (*domain.Room, error) {
	updatedRoom, err := m.softDelete(ctx, m.roomCollection, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ToDomainRoom(updatedRoom), nil
}

func (m *mongoRoomRepository) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error) {
	rooms, err := m.queryPaginated(
		ctx,
		m.roomCollection,
		skip, limit,
		bson.M{
			"isDelete": false,
		},
		bson.D{{Key: "createdAt", Value: 1}},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var results []domain.Room
	for _, room := range rooms {
		results = append(results, *ToDomainRoom(room))
	}

	return results, nil
}

func (m *mongoRoomRepository) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	room, err := m.getByID(ctx, m.roomCollection, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ToDomainRoom(room), nil
}

func (m *mongoRoomRepository) GetByFilter(
	ctx context.Context,
	ids []string,
	equipments []domain.Equipment, rules []domain.Rule,
	skip int, limit int,
) ([]domain.Room, error) {
	filter := bson.M{
		"isDelete": false,
	}

	if len(ids) > 0 {
		objIDs := make([]primitive.ObjectID, len(ids))
		for i, id := range ids {
			if objID, err := primitive.ObjectIDFromHex(id); err != nil {
				log.Printf("Invalid ID format: %v", err)
				return nil, err
			} else {
				objIDs[i] = objID
			}
		}
		filter["_id"] = bson.M{"$all": objIDs}
	}

	if len(equipments) > 0 {
		filter["equipments"] = bson.M{"$all": equipments}
	}
	if len(rules) > 0 {
		filter["rules"] = bson.M{"$all": rules}
	}

	rooms, err := m.queryPaginated(
		ctx,
		m.roomCollection,
		skip,
		limit,
		filter,
		bson.D{{Key: "name", Value: 1}},
	)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Room, len(rooms))
	for i, room := range rooms {
		res[i] = *ToDomainRoom(room)
	}

	return res, nil
}

func (m *mongoRoomRepository) QueryPaginatedAvailable(
	ctx context.Context,
	ids []string,
	equipments []domain.Equipment, rules []domain.Rule,
	startAt, endAt int64,
	skip int, limit int,
) ([]domain.Room, error) {
	filter := bson.M{
		"isDelete":     false,
		"reservations": bson.D{{Key: "$size", Value: 0}}, // in pipeline field
	}

	if len(ids) > 0 {
		objIDs := make([]primitive.ObjectID, len(ids))
		for i, id := range ids {
			if objID, err := primitive.ObjectIDFromHex(id); err != nil {
				log.Printf("Invalid ID format: %v", err)
				return nil, err
			} else {
				objIDs[i] = objID
			}
		}
		filter["_id"] = bson.M{"$all": objIDs}
	}

	if len(equipments) > 0 {
		filter["equipments"] = bson.M{"$all": equipments}
	}
	if len(rules) > 0 {
		filter["rules"] = bson.M{"$all": rules}
	}

	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "events"},
				{Key: "let", Value: bson.D{{Key: "roomId", Value: "$_id"}}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{
						{Key: "$match", Value: bson.D{
							{Key: "$expr", Value: bson.D{
								{Key: "$and", Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$roomReservation.roomID", "$$roomId"}}},
									bson.D{{Key: "$lt", Value: bson.A{"$startAt", endAt}}},
									bson.D{{Key: "$gt", Value: bson.A{"$endAt", startAt}}},
									bson.D{{Key: "$ne", Value: bson.A{"$roomReservation.reservationStatus", domain.ReservationStatus_Canceled}}},
									bson.D{{Key: "$eq", Value: bson.A{"$isDelete", false}}},
								}},
							}},
						}},
					},
				}},
				{Key: "as", Value: "reservations"},
			}},
		},
		{
			{Key: "$match", Value: filter},
		},
		{
			{Key: "$skip", Value: skip},
		},
		{
			{Key: "$limit", Value: limit},
		},
	}

	cursor, err := m.roomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []Room
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	res := make([]domain.Room, len(rooms))
	for i, room := range rooms {
		res[i] = *ToDomainRoom(&room)
	}

	return res, nil
}

func ToDomainRoom(room *Room) *domain.Room {
	domainRoom := domain.Room{
		ID:         common.ToPtr(room.ID.Hex()),
		Name:       room.Name,
		Capacity:   room.Capacity,
		Equipments: room.Equipments,
		Rules:      room.Rules,
		IsDelete:   room.IsDelete,
		CreatedAt:  room.CreatedAt,
		UpdatedAt:  room.UpdatedAt,
	}
	return &domainRoom
}
