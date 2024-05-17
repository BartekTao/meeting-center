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

type Event struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty"`
	Title           string               `bson:"title"`
	Description     *string              `bson:"description"`
	StartAt         int64                `bson:"startAt"`
	EndAt           int64                `bson:"endAt"`
	RoomReservation *RoomReservation     `bson:"roomReservation"`
	ParticipantsIDs []primitive.ObjectID `bson:"participantsIDs"`
	Notes           *string              `bson:"notes"`
	RemindAt        int64                `bson:"remindAt"`
	IsDelete        bool                 `bson:"isDelete"`
	Summary         string               `bson:"summary"`
	CreatedAt       int64                `bson:"createdAt"`
	CreatorID       primitive.ObjectID   `bson:"creatorID"`
	UpdatedAt       int64                `bson:"updatedAt"`
	UpdaterID       primitive.ObjectID   `bson:"updaterID"`
}

type RoomReservation struct {
	RoomID            primitive.ObjectID       `bson:"roomID"`
	ReservationStatus domain.ReservationStatus `bson:"reservationStatus"`
}

type mongoEventRepository struct {
	BaseRepository[Event]
	client          *mongo.Client
	eventCollection *mongo.Collection
}

func NewMongoEventRepository(client *mongo.Client) domain.EventRepository {
	return &mongoEventRepository{
		client:          client,
		eventCollection: client.Database("meetingCenter").Collection("events"),
	}
}

func (m *mongoEventRepository) Upsert(ctx context.Context, event domain.Event) (*domain.Event, error) {
	collection := m.eventCollection

	updaterID, err := primitive.ObjectIDFromHex(event.UpdaterID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, err
	}

	participantsIDs := make([]primitive.ObjectID, len(event.ParticipantsIDs))
	for i, participantsID := range event.ParticipantsIDs {
		objParticipantsID, err := primitive.ObjectIDFromHex(participantsID)
		if err != nil {
			log.Printf("Invalid ID format: %v", err)
			return nil, err
		}
		participantsIDs[i] = objParticipantsID
	}
	var roomReservation *RoomReservation
	if event.RoomReservation != nil {
		roomID, err := primitive.ObjectIDFromHex(*event.RoomReservation.RoomID)
		if err != nil {
			log.Printf("Invalid ID format: %v", err)
			return nil, err
		}
		roomReservation = &RoomReservation{
			RoomID:            roomID,
			ReservationStatus: event.RoomReservation.ReservationStatus,
		}
	}

	if event.ID == nil { // Insert new event
		currentTime := time.Now().Unix()

		newEvent := Event{
			Title:           event.Title,
			Description:     event.Description,
			StartAt:         event.StartAt,
			EndAt:           event.EndAt,
			ParticipantsIDs: participantsIDs,
			RoomReservation: roomReservation,
			Notes:           event.Notes,
			RemindAt:        event.RemindAt,
			IsDelete:        false,
			CreatedAt:       currentTime,
			CreatorID:       updaterID,
			UpdatedAt:       currentTime,
			UpdaterID:       updaterID,
		}

		result, err := collection.InsertOne(ctx, newEvent)
		if err != nil {
			log.Printf("Failed to insert new event: %v", err)
			return nil, err
		}
		newEvent.ID = result.InsertedID.(primitive.ObjectID)

		return ToDomainEvent(&newEvent), nil
	} else { // Update existing event
		id, err := primitive.ObjectIDFromHex(*event.ID)
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
				"title":           event.Title,
				"description":     event.Description,
				"startAt":         event.StartAt,
				"endAt":           event.EndAt,
				"roomReservation": roomReservation,
				"participantsIDs": participantsIDs,
				"notes":           event.Notes,
				"remindAt":        event.RemindAt,
				"updatedAt":       time.Now().Unix(),
				"updaterID":       updaterID,
			},
		}

		var updatedEvent Event
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

		return ToDomainEvent(&updatedEvent), nil
	}
}

func (m *mongoEventRepository) Delete(ctx context.Context, id string) (*domain.Event, error) {
	updatedEvent, err := m.softDelete(ctx, m.eventCollection, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ToDomainEvent(updatedEvent), nil
}

func (m *mongoEventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	event, err := m.getByID(ctx, m.eventCollection, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ToDomainEvent(event), nil
}

func (m *mongoEventRepository) GetByUsers(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error) {
	bsonIds, err := common.ToBsonIDs(ids)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participantsIDs": bson.M{"$in": bsonIds},
		"$or": []bson.M{
			{"startAt": bson.M{"$lt": endAt, "$gte": startAt}},
			{"endAt": bson.M{"$gt": startAt, "$lte": endAt}},
		},
		"isDelete": false,
	}
	events, err := m.BaseRepository.findAllByFilter(ctx, m.eventCollection, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	idSet := make(map[string]struct{})
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	result := make(map[string][]domain.Event)
	for _, event := range events {
		for _, participantID := range event.ParticipantsIDs {
			if _, exists := idSet[participantID.Hex()]; exists {
				result[participantID.Hex()] = append(result[participantID.Hex()], *ToDomainEvent(event))
			}
		}
	}

	return result, nil
}

func (m *mongoEventRepository) UpdateSummary(ctx context.Context, id string, summary string, updaterID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return false, err
	}
	filter := bson.M{
		"_id":      objID,
		"isDelete": false,
	}
	update := bson.M{
		"$set": bson.M{
			"summary":   summary,
			"updatedAt": time.Now().Unix(),
			"updaterID": updaterID,
		},
	}

	_, err = m.eventCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Document with the given ID not found or deleted: %v", err)
			return false, err
		} else {
			log.Printf("Failed to update: %v", err)
			return false, err
		}
	}
	return true, nil
}

func (m *mongoEventRepository) CheckAvailableRoom(ctx context.Context, roomID string, startAt, endAt int64) (bool, error) {
	objRoomID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return false, err
	}
	filter := bson.M{
		"roomReservation.roomID":            objRoomID,
		"roomReservation.reservationStatus": domain.ReservationStatus_Confirmed,
		"$or": []bson.M{
			{"startAt": bson.M{"$lt": endAt, "$gte": startAt}},
			{"endAt": bson.M{"$gt": startAt, "$lte": endAt}},
		},
		"isDelete": false,
	}
	res, err := m.findOneByFilter(ctx, m.eventCollection, filter)
	if err != nil {
		return false, err
	}
	return res == nil, nil
}

func (m *mongoEventRepository) GetAllWithRoomConfirmed(ctx context.Context, roomIDs []string, startAt, endAt int64) ([]domain.Event, error) {
	filter := bson.M{
		"roomReservation.reservationStatus": domain.ReservationStatus_Confirmed,
		"$or": []bson.M{
			{"startAt": bson.M{"$lte": endAt, "$gte": startAt}},
			{"endAt": bson.M{"$gte": startAt, "$lte": endAt}},
		},
		"isDelete": false,
	}
	if roomIDs != nil {
		objRoomIDs := make([]primitive.ObjectID, len(roomIDs))
		for i, roomID := range roomIDs {
			objRoomID, err := primitive.ObjectIDFromHex(roomID)
			if err != nil {
				log.Printf("Invalid ID format: %v", err)
				return nil, err
			}
			objRoomIDs[i] = objRoomID
		}
		filter["roomReservation.roomID"] = bson.M{"$in": objRoomIDs}
	} else {
		filter["roomReservation.roomID"] = bson.M{"$ne": nil}
	}
	events, err := m.findAllByFilter(ctx, m.eventCollection, filter)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Event, len(events))
	for i, event := range events {
		res[i] = *ToDomainEvent(event)
	}
	return res, nil
}

func ToDomainEvent(event *Event) *domain.Event {
	participantsIDs := make([]string, len(event.ParticipantsIDs))
	for i, participantsID := range event.ParticipantsIDs {
		participantsIDs[i] = participantsID.Hex()
	}
	domainRoom := domain.Event{
		ID:              common.ToPtr(event.ID.Hex()),
		Title:           event.Title,
		Description:     event.Description,
		StartAt:         event.StartAt,
		EndAt:           event.EndAt,
		RoomReservation: ToDomainRoomReservation(event.RoomReservation),
		ParticipantsIDs: participantsIDs,
		Notes:           event.Notes,
		RemindAt:        event.RemindAt,
		IsDelete:        event.IsDelete,
		CreatedAt:       event.CreatedAt,
		CreatorID:       event.CreatorID.Hex(),
		UpdatedAt:       event.UpdatedAt,
		UpdaterID:       event.UpdaterID.Hex(),
	}
	return &domainRoom
}

func ToDomainRoomReservation(roomReservation *RoomReservation) *domain.RoomReservation {
	if roomReservation == nil {
		return nil
	}
	domainRoom := domain.RoomReservation{
		RoomID:            common.ToPtr(roomReservation.RoomID.Hex()),
		ReservationStatus: roomReservation.ReservationStatus,
	}
	return &domainRoom
}

type MockEventRepository struct {
	mongoEventRepository
	client          *mongo.Client
	eventCollection *mongo.Collection
}

func (m *MockEventRepository) CheckAvailableRoom(ctx context.Context, roomID string, startAt, endAt int64) (bool, error) {
	return true, nil
}

func NewMockEventRepository(client *mongo.Client) domain.EventRepository {
	return &MockEventRepository{
		client:          client,
		eventCollection: client.Database("meetingCenter").Collection("events"),
	}
}
