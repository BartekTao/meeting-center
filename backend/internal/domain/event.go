package domain

import (
	"context"
)

type Event struct {
	ID              string   `json:"_id,omitempty"`
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	StartAt         int      `json:"startAt"`
	EndAt           int      `json:"endAt"`
	RoomID          *string  `json:"roomId"`
	ParticipantsIDs []string `json:"participantsIDs"`
	Notes           *string  `json:"notes"`
	RemindAt        int      `json:"remindAt"`
	IsDelete        bool     `json:"isDelete"`
	CreatedAt       int64    `json:"createdAt"`
	UpdatedAt       int64    `json:"updatedAt"`
	UpdaterId       string   `json:"updaterId"`
}

type EventRepository interface {
	UpsertEvent(ctx context.Context, event Event) (*Event, error)
	DeleteEvent(ctx context.Context, id string) (*Event, error)
}
