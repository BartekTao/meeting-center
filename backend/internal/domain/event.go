package domain

import (
	"context"
)

type ReservationStatus string

const (
	ReservationStatus_Canceled  ReservationStatus = "CANCELED"
	ReservationStatus_Confirmed ReservationStatus = "CONFIRMED"
)

type Event struct {
	ID              *string          `json:"_id,omitempty"`
	Title           string           `json:"title"`
	Description     *string          `json:"description"`
	StartAt         int64            `json:"startAt"`
	EndAt           int64            `json:"endAt"`
	RoomReservation *RoomReservation `json:"roomReservation"`
	ParticipantsIDs []string         `json:"participantsIDs"`
	Notes           *string          `json:"notes"`
	Summary         *string          `json:"summary"`
	RemindAt        int64            `json:"remindAt"`
	AttachedFile    File             `json:"attachedFile"`
	IsDelete        bool             `json:"isDelete"`
	CreatedAt       int64            `json:"createdAt"`
	CreatorID       string           `json:"creatorID"`
	UpdatedAt       int64            `json:"updatedAt"`
	UpdaterID       string           `json:"updaterID"`
}

type File struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type RoomReservation struct {
	RoomID            *string           `json:"roomID"`
	ReservationStatus ReservationStatus `json:"reservationStatus"`
}

type EventRepository interface {
	Upsert(ctx context.Context, event Event) (*Event, error)
	Delete(ctx context.Context, id string) (*Event, error)
	UpdateSummary(ctx context.Context, id string, summary string, updaterID string) (bool, error)
	GetByID(ctx context.Context, id string) (*Event, error)
	GetByUsers(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]Event, error)
	CheckAvailableRoom(ctx context.Context, roomID string, startAt, endAt int64) (bool, error)
	GetAllWithRoomConfirmed(ctx context.Context, roomIDs []string, startAt, endAt int64) ([]Event, error)
	GetRemindEvents(ctx context.Context, checkAt int64) ([]Event, error)
}
