package commands

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

type RoomHandler interface {
	UpsertRoom(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error)
	DeleteRoom(ctx context.Context, id string) (*domain.Room, error)
}

type roomHandler struct {
	roomRepository domain.RoomRepository
}

func NewRoomHandler(roomRepository domain.RoomRepository) RoomHandler {
	return roomHandler{roomRepository: roomRepository}
}

func (h roomHandler) UpsertRoom(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error) {
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

func (h roomHandler) DeleteRoom(ctx context.Context, id string) (*domain.Room, error) {
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
