package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/internal/common"
	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/model"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/middleware"
)

// BookedBy is the resolver for the bookedBy field.
func (r *bookingResolver) BookedBy(ctx context.Context, obj *model.Booking) (*domain.User, error) {
	return r.userService.GetByID(ctx, *obj.BookedBy.ID)
}

// UpsertRoom is the resolver for the upsertRoom field.
func (r *mutationResolver) UpsertRoom(ctx context.Context, room model.UpsertRoomInput) (*domain.Room, error) {
	claims, _ := ctx.Value(middleware.UserCtxKey).(*middleware.MeetingCenterClaims)

	upsertRoom := app.UpsertRoomRequest{
		ID:        room.ID,
		RoomID:    room.RoomID,
		Capacity:  room.Capacity,
		Equipment: room.Equipment,
		Rules:     room.Rules,
		UpdaterID: claims.Sub,
	}
	res, err := r.roomService.Upsert(ctx, upsertRoom)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteRoom is the resolver for the deleteRoom field.
func (r *mutationResolver) DeleteRoom(ctx context.Context, id string) (*domain.Room, error) {
	room, err := r.roomService.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// PaginatedRooms is the resolver for the paginatedRooms field.
func (r *queryResolver) PaginatedRooms(ctx context.Context, first *int, after *string) (*model.RoomConnection, error) {
	skip, err := common.DecodeCursor(after)
	if err != nil {
		return nil, err
	}

	rooms, err := r.roomService.QueryPaginated(ctx, *skip, *first)
	if err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, nil
	}

	edges := make([]*model.RoomEdge, len(rooms))
	for idx, room := range rooms {
		edges[idx] = &model.RoomEdge{
			Node:   &room,
			Cursor: common.EncodeCursor(idx + 1 + *skip),
		}
	}

	return &model.RoomConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &edges[0].Cursor,
			EndCursor:   &edges[len(edges)-1].Cursor,
		},
	}, nil
}

// Room is the resolver for the room field.
func (r *queryResolver) Room(ctx context.Context, id string) (*domain.Room, error) {
	room, err := r.roomService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// Bookings is the resolver for the bookings field.
func (r *roomResolver) Bookings(ctx context.Context, obj *domain.Room) ([]model.Booking, error) {
	panic(fmt.Errorf("not implemented: Bookings - bookings"))
}

// Booking returns graph.BookingResolver implementation.
func (r *Resolver) Booking() graph.BookingResolver { return &bookingResolver{r} }

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Room returns graph.RoomResolver implementation.
func (r *Resolver) Room() graph.RoomResolver { return &roomResolver{r} }

type (
	bookingResolver  struct{ *Resolver }
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	roomResolver     struct{ *Resolver }
)
