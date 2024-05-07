package resolvers

import (
	"github.com/BartekTao/nycu-meeting-room-api/internal/app/commands"
	"github.com/BartekTao/nycu-meeting-room-api/internal/meeting"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	meetingManager meeting.MeetingManager
	roomHandler    commands.UpsertRoomRequestHandler
}

func NewResolver(meetingManager meeting.MeetingManager, roomHandler commands.UpsertRoomRequestHandler) *Resolver {
	return &Resolver{
		meetingManager: meetingManager,
		roomHandler:    roomHandler,
	}
}
