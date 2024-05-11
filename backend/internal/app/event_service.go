package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type EventService interface {
	Upsert(ctx context.Context, req UpsertRoomRequest) (*domain.Event, error)
	Delete(ctx context.Context, id string) (*domain.Event, error)
	GetByID(ctx context.Context, id string) (*domain.Event, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Event, error)
}
