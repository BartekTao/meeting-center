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
	StartAt         int                `bson:"startAt"`
	EndAt           int                `bson:"endAt"`
	RoomID          *string            `bson:"roomId"`
	ParticipantsIDs []string           `bson:"participantsIDs"`
	Notes           *string            `bson:"notes"`
	RemindAt        int                `bson:"remindAt"`
	IsDelete        bool               `bson:"isDelete"`
	CreatedAt       int64              `bson:"createdAt"`
	UpdatedAt       int64              `bson:"updatedAt"`
	UpdaterId       string             `bson:"updaterId"`
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
			UpdatedAt:       currentTime,
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
		UpdatedAt:       event.UpdatedAt,
		UpdaterId:       event.UpdaterId,
	}
	return &domainRoom
}
