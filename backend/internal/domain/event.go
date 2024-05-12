package domain

import (
	"context"
)

type Event struct {
	ID              *string  `json:"_id,omitempty"`
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	StartAt         int64    `json:"startAt"`
	EndAt           int64    `json:"endAt"`
	RoomID          *string  `json:"roomId"`
	ParticipantsIDs []string `json:"participantsIDs"`
	Notes           *string  `json:"notes"`
	RemindAt        int64    `json:"remindAt"`
	IsDelete        bool     `json:"isDelete"`
	CreatedAt       int64    `json:"createdAt"`
	CreatorID       string   `json:"creatorID"`
	UpdatedAt       int64    `json:"updatedAt"`
	UpdaterID       string   `json:"updaterID"`
}

type EventRepository interface {
	Upsert(ctx context.Context, event Event) (*Event, error)
	Delete(ctx context.Context, id string) (*Event, error)
	GetByID(ctx context.Context, id string) (*Event, error)
	GetByUsers(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]Event, error)
}
