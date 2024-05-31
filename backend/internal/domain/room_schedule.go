package domain

import "context"

type RoomSchedule struct {
	Room      Room    `json:"room"`
	Schedules []Event `json:"schedules"`
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
