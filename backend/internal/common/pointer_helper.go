package common

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToPtr(s string) *string { return &s }

func ToBsonIDs(ids []string) ([]primitive.ObjectID, error) {
	bsonIds := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error converting ID to ObjectID:", err)
			return nil, err
		}
		bsonIds[i] = objID
	}
	return bsonIds, nil
}

func ToBsonID(id string) (*primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID:", err)
		return nil, err
	}
	return &objID, nil
}
