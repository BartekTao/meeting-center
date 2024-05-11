package resolvers

import (
	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/internal/meeting"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	meetingManager meeting.MeetingManager
	roomService    app.RoomService
}

func NewResolver(meetingManager meeting.MeetingManager, roomService app.RoomService) *Resolver {
	return &Resolver{
		meetingManager: meetingManager,
		roomService:    roomService,
	}
}
