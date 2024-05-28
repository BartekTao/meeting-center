package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type UpsertRoomRequest struct {
	ID         *string            `json:"_id,omitempty"`
	Name       string             `json:"name"`
	Capacity   int                `json:"capacity"`
	Equipments []domain.Equipment `json:"equipments"`
	Rules      []domain.Rule      `json:"rules"`
	UpdaterID  string             `json:"updaterID"`
}

type RoomService interface {
	Upsert(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error)
	Delete(ctx context.Context, id string) (*domain.Room, error)
	GetByID(ctx context.Context, id string) (*domain.Room, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error)
	QueryPaginatedRoomSchedule(
		ctx context.Context,
		roomIDs []string,
		equipments []domain.Equipment, rules []domain.Rule,
		startAt, endAt int64,
		skip int, limit int,
	) ([]domain.RoomSchedule, error)
	QueryPaginatedAvailable(
		ctx context.Context,
		ids []string,
		equipments []domain.Equipment, rules []domain.Rule,
		startAt, endAt int64,
		skip int, limit int,
	) ([]domain.Room, error)
}

type roomService struct {
	roomRepository   domain.RoomRepository
	roomScheduleRepo domain.RoomScheduleRepo
}

func NewRoomService(roomRepository domain.RoomRepository, roomScheduleRepo domain.RoomScheduleRepo) RoomService {
	return &roomService{roomRepository: roomRepository, roomScheduleRepo: roomScheduleRepo}
}

func (s *roomService) Upsert(ctx context.Context, req UpsertRoomRequest) (*domain.Room, error) {
	room := domain.Room{
		ID:         req.ID,
		Name:       req.Name,
		Capacity:   req.Capacity,
		Equipments: req.Equipments,
		Rules:      req.Rules,
		UpdaterID:  req.UpdaterID,
	}
	res, err := s.roomRepository.Upsert(ctx, room)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (s roomService) Delete(ctx context.Context, id string) (*domain.Room, error) {
	res, err := s.roomRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (s roomService) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	room, err := s.roomRepository.GetByID(ctx, id)
	return room, err
}

func (s roomService) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.Room, error) {
	rooms, err := s.roomRepository.QueryPaginated(ctx, skip, limit)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s roomService) QueryPaginatedRoomSchedule(
	ctx context.Context,
	roomIDs []string,
	equipments []domain.Equipment, rules []domain.Rule,
	startAt, endAt int64,
	skip int, limit int,
) ([]domain.RoomSchedule, error) {
	return s.roomScheduleRepo.QueryPaginated(
		ctx,
		roomIDs,
		equipments, rules,
		startAt, endAt,
		skip, limit,
	)
}

func (s roomService) QueryPaginatedAvailable(
	ctx context.Context,
	ids []string,
	equipments []domain.Equipment, rules []domain.Rule,
	startAt, endAt int64,
	skip int, limit int,
) ([]domain.Room, error) {
	return s.roomRepository.QueryPaginatedAvailable(
		ctx,
		ids,
		equipments, rules,
		startAt, endAt,
		skip, limit)
}
