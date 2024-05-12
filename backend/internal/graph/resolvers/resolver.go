package resolvers

import (
	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	roomService  app.RoomService
	eventService app.EventService
	userService  app.UserService
}

func NewResolver(roomService app.RoomService, eventService app.EventService, userService app.UserService) *Resolver {
	return &Resolver{
		roomService:  roomService,
		eventService: eventService,
		userService:  userService,
	}
}
