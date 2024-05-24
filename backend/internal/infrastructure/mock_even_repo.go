package infra

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockEventRepository struct {
	mongoEventRepository
	client          *mongo.Client
	eventCollection *mongo.Collection
	count           int
}

func (m *MockEventRepository) CheckAvailableRoom(ctx context.Context, roomID string, startAt, endAt int64) (bool, error) {
	if m.count == 0 {
		m.count += 1
		return true, nil
	} else {
		return false, nil
	}
}

func (m *MockEventRepository) Upsert(ctx context.Context, event domain.Event) (*domain.Event, error) {
	return &domain.Event{RoomReservation: &domain.RoomReservation{
		ReservationStatus: event.RoomReservation.ReservationStatus,
	}}, nil
}

func NewMockEventRepository(client *mongo.Client) domain.EventRepository {
	n := 0
	return &MockEventRepository{
		client:          client,
		eventCollection: client.Database("meetingCenter").Collection("events"),
		count:           n,
	}
}
