package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type UpsertEventRequest struct {
	ID              *string  `json:"_id,omitempty"`
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	StartAt         int      `json:"startAt"`
	EndAt           int      `json:"endAt"`
	RoomID          *string  `json:"roomId"`
	ParticipantsIDs []string `json:"participantsIDs"`
	Notes           *string  `json:"notes"`
	RemindAt        int      `json:"remindAt"`
	UpdaterId       string   `json:"updaterId"`
}

type EventService interface {
	Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error)
	Delete(ctx context.Context, id string) (*domain.Event, error)
	GetByID(ctx context.Context, id string) (*domain.Event, error)
}

type eventService struct {
	eventRepository domain.EventRepository
}

func NewEventService(eventRepository domain.EventRepository) EventService {
	return eventService{eventRepository: eventRepository}
}

func (h eventService) Upsert(ctx context.Context, req UpsertEventRequest) (*domain.Event, error) {
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

func (h eventService) Delete(ctx context.Context, id string) (*domain.Event, error) {
	res, err := h.eventRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h eventService) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	res, err := h.eventRepository.GetByID(ctx, id)
	return res, err
}

func (h eventService) GetParticipantEvents(ctx context.Context, ids []string, startAt, endAt int64) ([]*domain.Event, error) {
	h.eventRepository.GetByUsers(ctx, ids, startAt, endAt)
	return nil, nil
}
