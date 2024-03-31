package meeting

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"roomId"`
	Capacity  int                `bson:"capacity"`
	Equipment []string           `bson:"equipment"`
	Rules     []string           `bson:"rules"`
	CreatedAt int64              `bson:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt"`
	UpdaterId string             `bson:"updaterId"`
}
