package resolvers

import (
	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/internal/meeting"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	roomService  app.RoomService
	eventService app.EventService
}

func NewResolver(meetingManager meeting.MeetingManager, roomService app.RoomService, eventService app.EventService) *Resolver {
	return &Resolver{
		roomService:  roomService,
		eventService: eventService,
	}
}
