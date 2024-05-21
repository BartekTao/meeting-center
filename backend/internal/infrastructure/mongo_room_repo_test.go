package infra

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_mongoRoomRepository_QueryPaginatedAvailable(t *testing.T) {
	type args struct {
		ctx     context.Context
		startAt int64
		endAt   int64
		skip    int
		limit   int
	}

	var testMongoClient *mongo.Client
	var resource *dockertest.Resource

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			"MONGO_INITDB_DATABASE=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		t.Fatalf("Could not start MongoDB container: %s", err)
	}
	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		// Ping MongoDB to ensure it's ready
		err = client.Ping(ctx, nil)
		if err != nil {
			return err
		}
		testMongoClient = client
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to MongoDB container: %s", err)
	}
	db := testMongoClient.Database("testdb")
	repo := &mongoRoomRepository{roomCollection: db.Collection("rooms")}

	// insert 10 rooms into in-memory mongodb
	for i := 1; i <= 10; i++ {
		reservations := []bson.M{}
		if i%2 == 0 {
			reservations = append(reservations, bson.M{
				"roomReservation": bson.M{"roomID": primitive.NewObjectID().Hex(), "reservationStatus": domain.ReservationStatus_Confirmed},
				"startAt":         time.Now().Add(-1 * time.Hour).Unix(),
				"endAt":           time.Now().Add(1 * time.Hour).Unix(),
			})
		}
		room := bson.M{
			"_id":          primitive.NewObjectID(),
			"name":         fmt.Sprintf("room%d", i),
			"capacity":     10,
			"isDelete":     false,
			"reservations": reservations,
		}
		_, err := repo.roomCollection.InsertOne(context.TODO(), room)
		require.NoError(t, err)
	}
	// own design compare func
	compareRooms := func(got, want []domain.Room) bool {
		wantMap := make(map[string]int)
		for _, room := range want {
			wantMap[room.Name] = room.Capacity
		}

		for _, room := range got {
			if capacity, exists := wantMap[room.Name]; exists {
				if room.Capacity != capacity {
					return false
				}
				delete(wantMap, room.Name)
			}
		}
		return len(wantMap) == 0
	}
	//query available rooms in specific perioud
	tests := []struct {
		name    string
		m       *mongoRoomRepository
		args    args
		want    []domain.Room
		wantErr bool
	}{
		{
			name: "Available rooms within the time range",
			m:    repo,
			args: args{
				ctx:     context.Background(),
				startAt: time.Now().Add(-2 * time.Hour).Unix(),
				endAt:   time.Now().Add(2 * time.Hour).Unix(),
				skip:    0,
				limit:   10,
			},
			want: []domain.Room{
				{Name: "room1", Capacity: 10},
				{Name: "room3", Capacity: 10},
				{Name: "room5", Capacity: 10},
				{Name: "room7", Capacity: 10},
				{Name: "room9", Capacity: 10},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.QueryPaginatedAvailable(tt.args.ctx, tt.args.startAt, tt.args.endAt, tt.args.skip, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoRoomRepository.QueryPaginatedAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareRooms(got, tt.want) {
				t.Errorf("mongoRoomRepository.QueryPaginatedAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
