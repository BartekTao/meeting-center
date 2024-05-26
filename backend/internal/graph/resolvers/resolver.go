package resolvers

import (
	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/storage"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	roomService    app.RoomService
	eventService   app.EventService
	userService    app.UserService
	storageHandler storage.StorageHandler
}

func NewResolver(
	roomService app.RoomService,
	eventService app.EventService,
	userService app.UserService,
	storageHandler storage.StorageHandler,
) *Resolver {
	return &Resolver{
		roomService:    roomService,
		eventService:   eventService,
		userService:    userService,
		storageHandler: storageHandler,
	}
}
