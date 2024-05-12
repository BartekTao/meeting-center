package infra

import (
	"context"
	"log"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/common"
	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Sub           string             `bson:"sub"`
	Name          string             `bson:"name"`
	GivenName     string             `bson:"givenName"`
	FamilyName    string             `bson:"familyName"`
	Picture       string             `bson:"picture"`
	Locale        string             `bson:"locale"`
	Email         string             `bson:"email"`
	EmailVerified bool               `bson:"emailVerified"`
	CreatedAt     int64              `bson:"createdAt"`
	UpdatedAt     int64              `bson:"updatedAt"`
}

type mongoUserRepo struct {
	BaseRepository[User]
	client         *mongo.Client
	roomCollection *mongo.Collection
}

func NewMongoUserRepo(client *mongo.Client) domain.UserRepo {
	return &mongoUserRepo{
		client:         client,
		roomCollection: client.Database("meetingCenter").Collection("users"),
	}
}

func (r *mongoUserRepo) GetUser(ctx context.Context, sub string) (*domain.User, error) {
	filter := bson.M{"sub": sub}
	user, err := r.findOneByFilter(ctx, r.roomCollection, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ToDomainUser(user), nil
}

func (r *mongoUserRepo) SignUp(ctx context.Context, user domain.User) (*string, error) {
	currentTime := time.Now().Unix()
	newUser := User{
		Sub:           user.Sub,
		Name:          user.Name,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Picture:       user.Picture,
		Locale:        user.Locale,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	result, err := r.roomCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Failed to insert new user: %v", err)
		return nil, err
	}
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	return common.ToPtr(newUser.ID.Hex()), nil
}

func ToDomainUser(user *User) *domain.User {
	if user == nil {
		return nil
	}
	domainUser := domain.User{
		ID:            common.ToPtr(user.ID.Hex()),
		Sub:           user.Sub,
		Name:          user.Name,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Picture:       user.Picture,
		Locale:        user.Locale,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}
	return &domainUser
}
