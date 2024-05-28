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
	Room      Room    `bson:"room"`
	Schedules []Event `bson:"schedules"`
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
		"isDelete": false,
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

	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "events"},
				{Key: "let", Value: bson.D{{Key: "roomId", Value: "$_id"}}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{
						{Key: "$match", Value: bson.D{
							{Key: "$expr", Value: bson.D{
								{Key: "$and", Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$roomReservation.roomID", "$$roomId"}}},
									bson.D{{Key: "$lt", Value: bson.A{"$startAt", endAt}}},
									bson.D{{Key: "$gt", Value: bson.A{"$endAt", startAt}}},
									bson.D{{Key: "$ne", Value: bson.A{"$roomReservation.reservationStatus", domain.ReservationStatus_Canceled}}},
									bson.D{{Key: "$eq", Value: bson.A{"$isDelete", false}}},
								}},
							}},
						}},
					},
				}},
				{Key: "as", Value: "schedules"},
			}},
		},
		{
			{Key: "$match", Value: matchStage},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "room", Value: bson.D{
					{Key: "_id", Value: "$_id"},
					{Key: "name", Value: "$name"},
					{Key: "capacity", Value: "$capacity"},
					{Key: "equipments", Value: "$equipments"},
					{Key: "rules", Value: "$rules"},
					{Key: "isDelete", Value: "$isDelete"},
					{Key: "createdAt", Value: "$createdAt"},
					{Key: "creatorID", Value: "$creatorID"},
					{Key: "updatedAt", Value: "$updatedAt"},
					{Key: "updaterID", Value: "$updaterID"},
				}},
				{Key: "schedules", Value: 1},
			}},
		},
		{
			{Key: "$skip", Value: skip},
		},
		{
			{Key: "$limit", Value: limit},
		},
	}

	cursor, err := r.roomScheduleCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roomSchedules []RoomSchedule
	if err = cursor.All(ctx, &roomSchedules); err != nil {
		return nil, err
	}

	res := make([]domain.RoomSchedule, len(roomSchedules))
	for i, roomSchedule := range roomSchedules {
		res[i] = *ToDomainRoomSchedule(&roomSchedule)
	}

	return res, nil
}

func ToDomainRoomSchedule(roomSchedule *RoomSchedule) *domain.RoomSchedule {
	domainRoomSchedule := domain.RoomSchedule{
		Room:      *ToDomainRoom(&roomSchedule.Room),
		Schedules: make([]domain.Event, len(roomSchedule.Schedules)),
	}
	for i, schedule := range roomSchedule.Schedules {
		domainRoomSchedule.Schedules[i] = *ToDomainEvent(&schedule)
	}
	return &domainRoomSchedule
}
