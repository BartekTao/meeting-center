package meeting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MeetingManager interface {
	Reserve(ctx context.Context, startTime, endTime time.Time) error
	CancelReservation(ctx context.Context, id string) error
	// GetDetails() MeetingRoomDetails
}

type BasicMeetingManager struct {
	mongoClient *mongo.Client
}

func NewBasicMeetingManager(mongoClient *mongo.Client) *BasicMeetingManager {
	return &BasicMeetingManager{
		mongoClient: mongoClient,
	}
}

func (b *BasicMeetingManager) Reserve(ctx context.Context, startTime, endTime time.Time) error {
	return nil
}

func (b *BasicMeetingManager) CancelReservation(ctx context.Context, id string) error {
	return nil
}
