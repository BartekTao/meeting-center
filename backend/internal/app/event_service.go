package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
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
	UpdaterId       string   `json:"updaterId"`
}

type EventService interface {
	Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error)
	Delete(ctx context.Context, id string) (*domain.Event, error)
	GetByID(ctx context.Context, id string) (*domain.Event, error)
	GetUserEvents(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error)
}

type eventService struct {
	eventRepository domain.EventRepository
}

func NewEventService(eventRepository domain.EventRepository) EventService {
	return &eventService{eventRepository: eventRepository}
}

func (h *eventService) Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error) {
	event := domain.Event{
		ID:              req.ID,
		Title:           req.Title,
		Description:     req.Description,
		StartAt:         req.StartAt,
		EndAt:           req.EndAt,
		RoomID:          req.RoomID,
		ParticipantsIDs: req.ParticipantsIDs,
		Notes:           req.Notes,
		RemindAt:        req.RemindAt,
	}
	res, err := h.eventRepository.Upsert(ctx, event)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h *eventService) Delete(ctx context.Context, id string) (*domain.Event, error) {
	res, err := h.eventRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h *eventService) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	res, err := h.eventRepository.GetByID(ctx, id)
	return res, err
}

func (h *eventService) GetUserEvents(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error) {
	userEventMap, err := h.eventRepository.GetByUsers(ctx, ids, startAt, endAt)
	if err != nil {
		return nil, err
	}
	return userEventMap, nil
}

func (h *eventService) CheckUserAvailable(ctx context.Context, ids []string, startAt, endAt int64) (map[string]bool, error) {
	userEventMap, err := h.eventRepository.GetByUsers(ctx, ids, startAt, endAt)
	if err != nil {
		return nil, err
	}

	availableMap := make(map[string]bool)
	for userID, userEvent := range userEventMap {
		availableMap[userID] = userEvent != nil
	}
	return availableMap, nil
}
