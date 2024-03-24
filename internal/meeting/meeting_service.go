package meeting

import (
	"context"
	"time"
)

type MeetingManager interface {
	Reserve(ctx context.Context, startTime, endTime time.Time) error
	CancelReservation(id string) error
	// GetDetails() MeetingRoomDetails
}

type MeetingRepository interface {
}

type BasicMeetingManager struct {
}

func (b *BasicMeetingManager) Reserve(ctx context.Context, startTime, endTime time.Time) error {
	return nil
}

func (b *BasicMeetingManager) CancelReservation(id string) error {
	return nil
}
