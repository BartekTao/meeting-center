package meeting

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"roomID"`
	Capacity  int                `bson:"capacity"`
	Equipment []string           `bson:"equipment"`
	Rules     []string           `bson:"rules"`
	IsDelete  bool               `bson:"isDelete"`
	CreatedAt int64              `bson:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt"`
	UpdaterId string             `bson:"updaterId"`
}

type Event struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `bson:"title"`
	Description     *string            `bson:"description"`
	StartAt         int                `bson:"startAt"`
	EndAt           int                `bson:"endAt"`
	RoomID          *string            `bson:"roomId"`
	ParticipantsIDs []string           `bson:"participantsIDs"`
	Notes           *string            `bson:"notes"`
	RemindAt        int                `bson:"remindAt"`
	IsDelete        bool               `bson:"isDelete"`
	CreatedAt       int64              `bson:"createdAt"`
	UpdatedAt       int64              `bson:"updatedAt"`
	UpdaterId       string             `bson:"updaterId"`
}
