package app

import (
	"context"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.User, error)
}

type userService struct {
	userRepo domain.UserRepo
}

func NewUserService(userRepo domain.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (h *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	room, err := h.userRepo.GetByID(ctx, id)
	return room, err
}

func (h *userService) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.User, error) {
	rooms, err := h.userRepo.QueryPaginated(ctx, skip, limit)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
