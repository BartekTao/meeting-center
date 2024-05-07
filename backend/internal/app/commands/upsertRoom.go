package commands

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
	IsDelete  bool     `json:"isDelete"`
	UpdaterId string   `json:"updaterId"`
}

type UpsertRoomRequestHandler interface {
	Handle(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error)
}

type upsertRoomRequestHandler struct {
	roomRepository domain.RoomRepository
}

func NewUpsertRoomRequestHandler(roomRepository domain.RoomRepository) UpsertRoomRequestHandler {
	return upsertRoomRequestHandler{roomRepository: roomRepository}
}

func (h upsertRoomRequestHandler) Handle(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error) {
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
