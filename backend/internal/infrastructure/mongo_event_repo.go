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
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `bson:"title"`
	Description     *string            `bson:"description"`
	StartAt         int64              `bson:"startAt"`
	EndAt           int64              `bson:"endAt"`
	RoomID          *string            `bson:"roomId"`
	ParticipantsIDs []string           `bson:"participantsIDs"`
	Notes           *string            `bson:"notes"`
	RemindAt        int64              `bson:"remindAt"`
	IsDelete        bool               `bson:"isDelete"`
	Summary         string             `bson:"summary"`
	CreatedAt       int64              `bson:"createdAt"`
	CreatorID       string             `bson:"creatorID"`
	UpdatedAt       int64              `bson:"updatedAt"`
	UpdaterID       string             `bson:"updaterID"`
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

	if event.ID == nil { // Insert new event
		currentTime := time.Now().Unix()
		newEvent := Event{
			Title:           event.Title,
			Description:     event.Description,
			StartAt:         event.StartAt,
			EndAt:           event.EndAt,
			RoomID:          event.RoomID,
			ParticipantsIDs: event.ParticipantsIDs,
			Notes:           event.Notes,
			RemindAt:        event.RemindAt,
			IsDelete:        false,
			CreatedAt:       currentTime,
			CreatorID:       event.UpdaterID,
			UpdatedAt:       currentTime,
			UpdaterID:       event.UpdaterID,
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
				"roomId":          event.RoomID,
				"participantsIDs": event.ParticipantsIDs,
				"notes":           event.Notes,
				"remindAt":        event.RemindAt,
				"updatedAt":       time.Now().Unix(),
				"updaterID":       event.UpdaterID,
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
		"startAt":         bson.M{"$gte": startAt},
		"endAt":           bson.M{"$lte": endAt},
		"isDelete":        false,
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
			if _, exists := idSet[participantID]; exists {
				result[participantID] = append(result[participantID], *ToDomainEvent(event))
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

func ToDomainEvent(event *Event) *domain.Event {
	domainRoom := domain.Event{
		ID:              common.ToPtr(event.ID.Hex()),
		Title:           event.Title,
		Description:     event.Description,
		StartAt:         event.StartAt,
		EndAt:           event.EndAt,
		RoomID:          event.RoomID,
		ParticipantsIDs: event.ParticipantsIDs,
		Notes:           event.Notes,
		RemindAt:        event.RemindAt,
		IsDelete:        event.IsDelete,
		CreatedAt:       event.CreatedAt,
		CreatorID:       event.CreatorID,
		UpdatedAt:       event.UpdatedAt,
		UpdaterID:       event.UpdaterID,
	}
	return &domainRoom
}
