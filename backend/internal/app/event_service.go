package app

import (
	"context"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/lock"
)

type UpsertEventRequest struct {
	ID              *string  `json:"_id,omitempty"`
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	StartAt         int64    `json:"startAt"`
	EndAt           int64    `json:"endAt"`
	RoomID          *string  `json:"roomId"`
	ParticipantsIDs []string `json:"participantsIDs"`
	Notes           *string  `json:"notes"`
	RemindAt        int64    `json:"remindAt"`
	UpdaterID       string   `json:"updaterID"`
}

type EventService interface {
	Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error)
	Delete(ctx context.Context, id string) (*domain.Event, error)
	GetByID(ctx context.Context, id string) (*domain.Event, error)
	GetUserEvents(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error)
	UpdateSummary(ctx context.Context, id string, summary string, updaterID string) (bool, error)
}

type eventService struct {
	eventRepository domain.EventRepository
	roomRepository  domain.RoomRepository
	locker          lock.DistributedLocker
}

func NewEventService(eventRepository domain.EventRepository, locker lock.DistributedLocker) EventService {
	return &eventService{
		eventRepository: eventRepository,
		locker:          locker,
	}
}

func (s *eventService) Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error) {
	event := domain.Event{
		ID:              req.ID,
		Title:           req.Title,
		Description:     req.Description,
		StartAt:         req.StartAt,
		EndAt:           req.EndAt,
		ParticipantsIDs: req.ParticipantsIDs,
		Notes:           req.Notes,
		RemindAt:        req.RemindAt,
		UpdaterID:       req.UpdaterID,
	}

	if req.RoomID == nil {
		res, err := s.eventRepository.Upsert(ctx, event)
		return res, err
	} else {
		event.RoomReservation = &domain.RoomReservation{
			RoomID:            req.RoomID,
			ReservationStatus: domain.ReservationStatus_Confirmed,
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	key := *req.RoomID
	locked, err := s.locker.TryLockWithWait(key, 500*time.Millisecond, 3)
	if err != nil {
		return nil, err
	}
	defer s.locker.Unlock(locked)

	available, err := s.eventRepository.CheckAvailableRoom(ctx, *event.RoomReservation.RoomID, event.StartAt, event.EndAt)
	if err != nil {
		return nil, err
	}

	if !available {
		event.RoomReservation.ReservationStatus = domain.ReservationStatus_Canceled
	}

	res, err := s.eventRepository.Upsert(ctx, event)

	return res, err
}

func (s *eventService) Delete(ctx context.Context, id string) (*domain.Event, error) {
	res, err := s.eventRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (s *eventService) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	res, err := s.eventRepository.GetByID(ctx, id)
	return res, err
}

func (s *eventService) GetUserEvents(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error) {
	userEventMap, err := s.eventRepository.GetByUsers(ctx, ids, startAt, endAt)
	if err != nil {
		return nil, err
	}
	return userEventMap, nil
}

func (s *eventService) CheckUserAvailable(ctx context.Context, ids []string, startAt, endAt int64) (map[string]bool, error) {
	userEventMap, err := s.eventRepository.GetByUsers(ctx, ids, startAt, endAt)
	if err != nil {
		return nil, err
	}

	availableMap := make(map[string]bool)
	for userID, userEvent := range userEventMap {
		availableMap[userID] = userEvent != nil
	}
	return availableMap, nil
}

func (s *eventService) UpdateSummary(ctx context.Context, id string, summary string, updaterID string) (bool, error) {
	return s.eventRepository.UpdateSummary(ctx, id, summary, updaterID)
}

func (s *eventService) PaginatedAvailableRooms(
	ctx context.Context,
	startAt, endAt int64,
	equipments []domain.Equipment, rules []domain.Rule,
	skip, limit int,
) ([]domain.Event, error) {
	rooms, err := s.roomRepository.GetAll(ctx, equipments, rules)
	if err != nil {
		return nil, err
	}
	roomIDs := make([]string, len(rooms))
	for i, room := range rooms {
		roomIDs[i] = *room.ID
	}
	events, err := s.eventRepository.GetAll(ctx, roomIDs, startAt, endAt)
	return events, err
}

func (s *eventService) PaginatedRoomStatus(
	ctx context.Context,
	startAt, endAt int64,
	roomIDs []string,
	skip, limit int,
) ([]domain.Event, error) {
	events, err := s.eventRepository.GetAll(ctx, roomIDs, startAt, endAt)
	return events, err
}
