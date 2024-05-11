package app

import (
	"context"
	"errors"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type UpsertRoomRequest struct {
	ID        *string  `json:"_id,omitempty"`
	RoomID    string   `json:"roomID"`
	Capacity  int      `json:"capacity"`
	Equipment []string `json:"equipment"`
	Rules     []string `json:"rules"`
	IsDelete  bool     `json:"isDelete"`
	UpdaterId string   `json:"updaterId"`
}

type RoomService interface {
	UpsertRoom(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error)
	DeleteRoom(ctx context.Context, id string) (*domain.Room, error)
	GetByID(ctx context.Context, id string) (*domain.Room, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error)
}

type roomService struct {
	roomRepository domain.RoomRepository
}

func NewRoomService(roomRepository domain.RoomRepository) RoomService {
	return roomService{roomRepository: roomRepository}
}

func (h roomService) UpsertRoom(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error) {
	room := domain.Room{
		ID:        req.ID,
		RoomID:    req.RoomID,
		Capacity:  req.Capacity,
		Equipment: req.Equipment,
		Rules:     req.Rules,
		IsDelete:  req.IsDelete,
		UpdaterId: req.UpdaterId,
	}
	res, err := h.roomRepository.UpsertRoom(ctx, room)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h roomService) DeleteRoom(ctx context.Context, id string) (*domain.Room, error) {
	room, err := h.roomRepository.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, errors.New("room not found")
	}

	room.DeleteRoom()
	res, err := h.roomRepository.UpsertRoom(ctx, *room)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h roomService) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	room, err := h.roomRepository.GetRoomByID(ctx, id)
	return room, err
}

func (h roomService) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error) {
	rooms, err := h.roomRepository.QueryPaginatedRoom(ctx, skip, limit)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
