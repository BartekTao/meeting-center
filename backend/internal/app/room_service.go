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

type QueryPaginatedRoomScheduleResult struct {
	Room      domain.Room `json:"room"`
	Schedules []Schedule  `json:"schedules"`
}

type Schedule struct {
	StartAt int64 `json:"startAt"`
	EndAt   int64 `json:"endAt"`
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
	) ([]QueryPaginatedRoomScheduleResult, error)
}

type roomService struct {
	roomRepository domain.RoomRepository
	eventRepo      domain.EventRepository
}

func NewRoomService(roomRepository domain.RoomRepository, eventRepo domain.EventRepository) RoomService {
	return &roomService{roomRepository: roomRepository, eventRepo: eventRepo}
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
) ([]QueryPaginatedRoomScheduleResult, error) {
	var rooms []domain.Room
	var err error
	if len(roomIDs) == 0 {
		rooms, err = s.roomRepository.GetByFilter(ctx, nil, equipments, rules, skip, limit)
		if err != nil {
			return nil, err
		}
		roomIDs = make([]string, len(rooms))
		for i, room := range rooms {
			roomIDs[i] = *room.ID
		}
	} else {
		rooms, err = s.roomRepository.GetByFilter(ctx, roomIDs, nil, nil, skip, limit)
		if err != nil {
			return nil, err
		}
	}

	roomMap := make(map[string]domain.Room)
	for _, room := range rooms {
		roomMap[*room.ID] = room
	}

	events, err := s.eventRepo.GetAllWithRoomConfirmed(ctx, roomIDs, startAt, endAt)
	if err != nil {
		return nil, err
	}

	scheduleMap := make(map[string][]Schedule)
	for _, event := range events {
		reservationRoomID := *event.RoomReservation.RoomID
		schedule := Schedule{StartAt: event.StartAt, EndAt: event.EndAt}
		if schedules, ok := scheduleMap[reservationRoomID]; ok {
			scheduleMap[reservationRoomID] = append(schedules, schedule)
		} else {
			scheduleMap[reservationRoomID] = []Schedule{schedule}
		}
	}

	res := make([]QueryPaginatedRoomScheduleResult, len(rooms))
	for i, room := range rooms {
		if schedules, ok := scheduleMap[*room.ID]; ok {
			res[i] = QueryPaginatedRoomScheduleResult{
				Room:      room,
				Schedules: schedules,
			}
		} else {
			res[i] = QueryPaginatedRoomScheduleResult{
				Room: room,
			}
		}
	}

	return res, nil
}
