package domain

import (
	"context"
)

type Room struct {
	ID        *string  `json:"_id,omitempty"`
	RoomID    string   `json:"roomID"`
	Capacity  int      `json:"capacity"`
	Equipment []string `json:"equipment"`
	Rules     []string `json:"rules"`
	IsDelete  bool     `json:"isDelete"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
	UpdaterId string   `json:"updaterId"`
}

func (r *Room) DeleteRoom() {
	r.IsDelete = true
}

type RoomRepository interface {
	UpsertRoom(ctx context.Context, room Room) (*Room, error)
	DeleteRoom(ctx context.Context, id string) (*Room, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	QueryPaginatedRoom(ctx context.Context, skip int, limit int) ([]Room, error)
}
