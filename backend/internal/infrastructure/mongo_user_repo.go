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
	userCollection *mongo.Collection
}

func NewMongoUserRepo(client *mongo.Client) domain.UserRepo {
	return &mongoUserRepo{
		client:         client,
		userCollection: client.Database("meetingCenter").Collection("users"),
	}
}

func (r *mongoUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.getByID(ctx, r.userCollection, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ToDomainUser(user), nil
}

func (r *mongoUserRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	bsonIds, err := common.ToBsonIDs(ids)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": bson.M{"$in": bsonIds},
	}
	users, err := r.findAllByFilter(ctx, r.userCollection, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = *ToDomainUser(user)
	}
	return domainUsers, nil
}

func (r *mongoUserRepo) GetUserBySub(ctx context.Context, sub string) (*domain.User, error) {
	filter := bson.M{"sub": sub}
	user, err := r.findOneByFilter(ctx, r.userCollection, filter)
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
	result, err := r.userCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Failed to insert new user: %v", err)
		return nil, err
	}
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	return common.ToPtr(newUser.ID.Hex()), nil
}

func (r *mongoUserRepo) QueryPaginated(ctx context.Context, skip int, limit int) ([]domain.User, error) {
	users, err := r.queryPaginated(
		ctx,
		r.userCollection,
		skip, limit, bson.M{},
		bson.D{{Key: "CreatedAt", Value: 1}},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var results []domain.User
	for _, room := range users {
		results = append(results, *ToDomainUser(room))
	}

	return results, nil
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
