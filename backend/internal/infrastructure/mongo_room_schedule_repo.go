package infra

import (
	"context"
	"log"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomSchedule struct {
	Room      Room       `bson:"room"`
	Schedules []Schedule `bson:"schedules"`
}

type Schedule struct {
	StartAt int64 `bson:"startAt"`
	EndAt   int64 `bson:"endAt"`
}

type mongoRoomScheduleRepository struct {
	BaseRepository[RoomSchedule]
	client                 *mongo.Client
	roomScheduleCollection *mongo.Collection
}

func NewRoomScheduleRepository(client *mongo.Client) domain.RoomScheduleRepo {
	return &mongoRoomScheduleRepository{
		client:                 client,
		roomScheduleCollection: client.Database("meetingCenter").Collection("rooms"),
	}
}

func (r *mongoRoomScheduleRepository) QueryPaginated(
	ctx context.Context,
	roomIDs []string,
	equipments []domain.Equipment, rules []domain.Rule,
	startAt, endAt int64,
	skip int, limit int,
) ([]domain.RoomSchedule, error) {
	matchStage := bson.M{
		"schedules.roomReservation.reservationStatus": domain.ReservationStatus_Confirmed,
		"$or": []bson.M{
			{"schedules.startAt": bson.M{"$lte": endAt, "$gte": startAt}},
			{"schedules.endAt": bson.M{"$lte": endAt, "$gte": startAt}},
		},
	}

	if len(roomIDs) > 0 {
		objIDs := make([]primitive.ObjectID, len(roomIDs))
		for i, id := range roomIDs {
			if objID, err := primitive.ObjectIDFromHex(id); err != nil {
				log.Printf("Invalid ID format: %v", err)
				return nil, err
			} else {
				objIDs[i] = objID
			}
		}
		matchStage["_id"] = bson.M{"$in": objIDs}
	}

	if len(equipments) > 0 {
		matchStage["equipments"] = bson.M{"$all": equipments}
	}
	if len(rules) > 0 {
		matchStage["rules"] = bson.M{"$all": rules}
	}

	// Lookup stage
	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":         "events",
			"localField":   "_id",
			"foreignField": "roomReservation.roomID",
			"as":           "schedules",
		},
	}

	// Unwind stage
	unwindStage := bson.M{
		"$unwind": bson.M{
			"path":                       "$schedules",
			"preserveNullAndEmptyArrays": true,
		},
	}

	// Group stage
	groupStage := bson.M{
		"$group": bson.M{
			"_id": "$_id",
			"room": bson.M{
				"$first": "$$ROOT",
			},
			"schedules": bson.M{
				"$push": bson.M{
					"startAt": "$schedules.startAt",
					"endAt":   "$schedules.endAt",
				},
			},
		},
	}

	// Project stage
	projectStage := bson.M{
		"$project": bson.M{
			"_id": 0,
			"room": bson.M{
				"_id":        "$room._id",
				"name":       "$room.name",
				"capacity":   "$room.capacity",
				"equipments": "$room.equipments",
				"rules":      "$room.rules",
				"isDelete":   "$room.isDelete",
				"createdAt":  "$room.createdAt",
				"creatorID":  "$room.creatorID",
				"updatedAt":  "$room.updatedAt",
				"updaterID":  "$room.updaterID",
			},
			"schedules": 1,
		},
	}

	// Sort stage
	sortStage := bson.M{
		"$sort": bson.M{"room.name": 1},
	}

	// Pagination stages
	skipStage := bson.M{
		"$skip": skip,
	}

	limitStage := bson.M{
		"$limit": limit,
	}

	// Build pipeline
	pipeline := []bson.M{
		lookupStage,
		unwindStage,
		{"$match": matchStage},
		groupStage,
		projectStage,
		sortStage,
		skipStage,
		limitStage,
	}

	cursor, err := r.roomScheduleCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []RoomSchedule
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	domainRoomSchedules := make([]domain.RoomSchedule, len(results))
	for i, domainRoomSchedule := range results {
		domainRoomSchedules[i] = *ToDomainRoomSchedule(&domainRoomSchedule)
	}

	return domainRoomSchedules, nil
}

func ToDomainRoomSchedule(roomSchedule *RoomSchedule) *domain.RoomSchedule {
	domainRoomSchedule := domain.RoomSchedule{
		Room:      *ToDomainRoom(&roomSchedule.Room),
		Schedules: make([]domain.Schedule, len(roomSchedule.Schedules)),
	}
	for i, schedule := range roomSchedule.Schedules {
		domainRoomSchedule.Schedules[i] = domain.Schedule{
			StartAt: schedule.StartAt,
			EndAt:   schedule.EndAt,
		}
	}
	return &domainRoomSchedule
}
