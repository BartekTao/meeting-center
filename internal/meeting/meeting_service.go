package meeting

import (
	"context"
	"fmt"

	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
)

type MeetingManager interface {
	UpsertRoom(ctx context.Context, room model.UpsertRoomInput) (*Room, error)
	DeleteRoom(ctx context.Context, id string) (*Room, error)
	QueryPaginatedRoom(ctx context.Context, first int, after string) (*Room, error)
	QueryRoom(ctx context.Context, id string) (*Room, error)
	QueryRooms(ctx context.Context, id string) (*Room, error)
}

type MeetingRepository interface {
	UpsertRoom(ctx context.Context, upsertRoomInput model.UpsertRoomInput) (*Room, error)
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
	name := "try"
	fmt.Println(name)
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}
func (b *BasicMeetingManager) QueryPaginatedRoom(ctx context.Context, first int, after string) (*Room, error) {
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}
func (b *BasicMeetingManager) QueryRoom(ctx context.Context, id string) (*Room, error) {
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}
func (b *BasicMeetingManager) QueryRooms(ctx context.Context, id string) (*Room, error) {
	panic(fmt.Errorf("not implemented: DeleteRoom - deleteRoom"))
}
