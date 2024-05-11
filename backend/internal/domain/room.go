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
	Upsert(ctx context.Context, room Room) (*Room, error)
	Delete(ctx context.Context, id string) (*Room, error)
	GetByID(ctx context.Context, id string) (*Room, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]Room, error)
}
