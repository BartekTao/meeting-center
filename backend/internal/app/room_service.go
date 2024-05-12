package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type UpsertRoomRequest struct {
	ID        *string  `json:"_id,omitempty"`
	RoomID    string   `json:"roomID"`
	Capacity  int      `json:"capacity"`
	Equipment []string `json:"equipment"`
	Rules     []string `json:"rules"`
	UpdaterID string   `json:"updaterID"`
}

type RoomService interface {
	Upsert(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error)
	Delete(ctx context.Context, id string) (*domain.Room, error)
	GetByID(ctx context.Context, id string) (*domain.Room, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error)
}

type roomService struct {
	roomRepository domain.RoomRepository
}

func NewRoomService(roomRepository domain.RoomRepository) RoomService {
	return &roomService{roomRepository: roomRepository}
}

func (h *roomService) Upsert(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error) {
	room := domain.Room{
		ID:        req.ID,
		RoomID:    req.RoomID,
		Capacity:  req.Capacity,
		Equipment: req.Equipment,
		Rules:     req.Rules,
		UpdaterID: req.UpdaterID,
	}
	res, err := h.roomRepository.Upsert(ctx, room)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h roomService) Delete(ctx context.Context, id string) (*domain.Room, error) {
	res, err := h.roomRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h roomService) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	room, err := h.roomRepository.GetByID(ctx, id)
	return room, err
}

func (h roomService) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error) {
	rooms, err := h.roomRepository.QueryPaginated(ctx, skip, limit)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
