package meeting

import (
	"context"
	"fmt"

	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
)

type MeetingManager interface {
	//UpsertRoom(ctx context.Context, room model.UpsertRoomInput) (*Room, error)
	UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*Room, error)
	DeleteRoom(ctx context.Context, id string) (*Room, error)
	QueryPaginatedRoom(ctx context.Context, first int, last int, before string, after string) (*model.RoomConnection, error)
	QueryRoom(ctx context.Context, id string) (*Room, error)
	QueryRooms(ctx context.Context, id string) (*Room, error)

	UpsertEvent(ctx context.Context, upsertEventInput model.UpsertEventInput) (*Event, error)
	DeleteEvent(ctx context.Context, id string) (*Event, error)
}

type MeetingRepository interface {
	UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*Room, error)
	DeleteRoom(ctx context.Context, id string) (*Room, error)
	QueryPaginatedRoom(ctx context.Context, first int, last int, before string, after string) (*model.RoomConnection, error)

	UpsertEvent(ctx context.Context, upsertEventInput model.UpsertEventInput) (*Event, error)
	DeleteEvent(ctx context.Context, id string) (*Event, error)
}

type BasicMeetingManager struct {
	meetingRepository MeetingRepository
}

func NewBasicMeetingManager(meetingRepository MeetingRepository) *BasicMeetingManager {
	return &BasicMeetingManager{
		meetingRepository: meetingRepository,
	}
}

func (b *BasicMeetingManager) UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*Room, error) {
	room, err := b.meetingRepository.UpsertRoom(ctx, upsertRoomInput)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (b *BasicMeetingManager) DeleteRoom(ctx context.Context, id string) (*Room, error) {
	room, err := b.meetingRepository.DeleteRoom(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (b *BasicMeetingManager) QueryPaginatedRoom(ctx context.Context, first int, last int, before string, after string) (*model.RoomConnection, error) {
	// panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
	room, err := b.meetingRepository.QueryPaginatedRoom(ctx, first, last, before, after)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (b *BasicMeetingManager) QueryRoom(ctx context.Context, id string) (*Room, error) {
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}

func (b *BasicMeetingManager) QueryRooms(ctx context.Context, id string) (*Room, error) {
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}

func (b *BasicMeetingManager) UpsertEvent(ctx context.Context, upsertEventInput model.UpsertEventInput) (*Event, error) {
	event, err := b.meetingRepository.UpsertEvent(ctx, upsertEventInput)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (b *BasicMeetingManager) DeleteEvent(ctx context.Context, id string) (*Event, error) {
	event, err := b.meetingRepository.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}
	return event, nil
}
