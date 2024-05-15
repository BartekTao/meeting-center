package domain

import "context"

type RoomSchedule struct {
	Room      Room       `json:"room"`
	Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
	StartAt int64 `json:"startAt"`
	EndAt   int64 `json:"endAt"`
}

type RoomScheduleRepo interface {
	QueryPaginated(
		ctx context.Context,
		roomIDs []string,
		equipments []Equipment, rules []Rule,
		startAt, endAt int64,
		skip int, limit int,
	) ([]RoomSchedule, error)
}
