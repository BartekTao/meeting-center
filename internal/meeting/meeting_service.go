package meeting

import (
	"context"
	"time"
)

type MeetingManager interface {
	Reserve(ctx context.Context, startTime, endTime time.Time) error
	CancelReservation(ctx context.Context, id string) error
	// GetDetails() MeetingRoomDetails
}

type MeetingRepository interface {
	// TODO: write some database interface method
}

type basicMeetingManager struct {
	meetingRepository MeetingRepository
}

func NewBasicMeetingManager(meetingRepository MeetingRepository) *basicMeetingManager {
	return &basicMeetingManager{
		meetingRepository: meetingRepository,
	}
}

func (b *basicMeetingManager) Reserve(ctx context.Context, startTime, endTime time.Time) error {
	return nil
}

func (b *basicMeetingManager) CancelReservation(ctx context.Context, id string) error {
	return nil
}
