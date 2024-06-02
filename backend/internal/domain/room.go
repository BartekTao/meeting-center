package domain

import (
	"context"
)

type Rule string

const (
	RULE_NO_FOOD  Rule = "NO_FOOD"
	RULE_NO_DRINK Rule = "NO_DRINK"
)

type Equipment string

const (
	EQUIPMENT_PROJECTOR  Equipment = "PROJECTOR"
	EQUIPMENT_TABLE      Equipment = "TABLE"
	EQUIPMENT_TV         Equipment = "TV"
	EQUIPMENT_CAMERA     Equipment = "CAMERA"
	EQUIPMENT_WHITEBOARD Equipment = "WHITEBOARD"
)

type Room struct {
	ID         *string     `json:"_id,omitempty"`
	Name       string      `json:"name"`
	Capacity   int         `json:"capacity"`
	Equipments []Equipment `json:"equipments"`
	Rules      []Rule      `json:"rules"`
	IsDelete   bool        `json:"isDelete"`
	CreatedAt  int64       `json:"createdAt"`
	CreatorID  string      `json:"creatorID"`
	UpdatedAt  int64       `json:"updatedAt"`
	UpdaterID  string      `json:"updaterID"`
}

type RoomRepository interface {
	Upsert(ctx context.Context, room Room) (*Room, error)
	Delete(ctx context.Context, id string) (*Room, error)
	GetByID(ctx context.Context, id string) (*Room, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]Room, error)
	QueryPaginatedAvailable(
		ctx context.Context,
		ids []string,
		equipments []Equipment, rules []Rule,
		startAt, endAt int64,
		skip int, limit int) ([]Room, error)
	GetByFilter(
		ctx context.Context,
		ids []string,
		equipments []Equipment, rules []Rule,
		skip int, limit int,
	) ([]Room, error)
}
