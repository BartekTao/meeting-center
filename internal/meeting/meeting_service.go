package meeting

import (
	"context"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
)

type MeetingManager interface {
	Reserve(ctx context.Context, startTime, endTime time.Time) error
	CancelReservation(ctx context.Context, id string) error
	// GetDetails() MeetingRoomDetails
}

type MeetingRepository interface {
	CreateRoom(ctx context.Context, createRoomInput model.CreateRoomInput) (Room, error)
}

type BasicMeetingManager struct {
	meetingRepository MeetingRepository
}

func NewBasicMeetingManager(meetingRepository MeetingRepository) *BasicMeetingManager {
	return &BasicMeetingManager{
		meetingRepository: meetingRepository,
	}
}

func (b *BasicMeetingManager) Reserve(ctx context.Context, startTime, endTime time.Time) error {
	return nil
}

func (b *BasicMeetingManager) CancelReservation(ctx context.Context, id string) error {
	return nil
}
