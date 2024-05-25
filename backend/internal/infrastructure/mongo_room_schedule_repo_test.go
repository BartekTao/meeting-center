package infra

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_mongoRoomScheduleRepository_QueryPaginated(t *testing.T) {
	type args struct {
		ctx        context.Context
		roomIDs    []string
		equipments []domain.Equipment
		rules      []domain.Rule
		startAt    int64
		endAt      int64
		skip       int
		limit      int
	}

	var testMongoClient *mongo.Client
	var pool *dockertest.Pool
	var resource *dockertest.Resource

	///////////////// Set up in-memory mongodb ////////////////////////////////

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

	// Wait for MongoDB to start up
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

	repo := NewRoomScheduleRepository(testMongoClient)
	testRoomScheduleRepo, ok := repo.(*mongoRoomScheduleRepository)
	if !ok {
		log.Printf("Failed to type assert repo to mongoEventRepository")
		t.Skip("Skipping test due to type assertion failure")
		return
	}

	////////////////// Setup room reservation mongodb data ////////////////////////////////////////////////

	for i := 0; i <= 2; i++ {
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
			"capacity":     20,
			"isDelete":     false,
			"reservations": reservations,
		}
		schedules := []Schedule{
			bson.M{
				"startAt": 100,
				"endAt": 200,
			},
			bson.M{
				"startAt": 400,
				"endAt": 500,
			},
			bson.M{
				"startAt": 700,
				"endAt": 800,
			},
		}
		roomSchedule := bson.M{
			"room": room,
			"schedules": schedules, 
		}
		_, err := testRoomScheduleRepo.roomScheduleCollection.InsertOne(context.TODO(), roomSchedule)
		require.NoError(t, err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	tests := []struct {
		name    string
		r       *mongoRoomScheduleRepository
		args    args
		want    []domain.RoomSchedule
		wantErr bool
	}{
		{
			name: "Successful Upsert",
			r:    testRoomScheduleRepo, // Initialize with appropriate values
			args: args{
				ctx:   context.Background(), // Use context appropriate for testing
				roomIDs: ,
				equipments: ,
				rules: , 
				startAt: ,
				endAt: ,
				skip: ,
				limit: ,
			},
			want:    domain.ReservationStatus_Confirmed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.QueryPaginated(tt.args.ctx, tt.args.roomIDs, tt.args.equipments, tt.args.rules, tt.args.startAt, tt.args.endAt, tt.args.skip, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoRoomScheduleRepository.QueryPaginated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoRoomScheduleRepository.QueryPaginated() = %v, want %v", got, tt.want)
			}
		})
	}
}
