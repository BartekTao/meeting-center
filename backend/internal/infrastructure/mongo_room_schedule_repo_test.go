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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	eventRepo := NewMongoEventRepository(testMongoClient)
	testEventRepo, ok := eventRepo.(*mongoEventRepository)
	if !ok {
		log.Printf("Failed to type assert repo to mongoEventRepository")
		t.Skip("Skipping test due to type assertion failure")
		return
	}

	////////////////// Setup room reservation mongodb data ////////////////////////////////////////////////

	testIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID()}
	testStr := "test"
	for i := 0; i <= 2; i++ {

		room := bson.M{
			"_id":        testIDs[i],
			"name":       fmt.Sprintf("room%d", i),
			"capacity":   20,
			"equipments": []domain.Equipment{domain.EQUIPMENT_TABLE, domain.EQUIPMENT_CAMERA},
			"rules":      []domain.Rule{domain.RULE_NO_FOOD, domain.RULE_NO_DRINK},
			"isDelete":   false,
		}
		var schedules Event
		if i%2 == 0 {
			schedules = Event{
				ID:              primitive.NewObjectID(),
				Title:           testStr,
				Description:     &testStr,
				StartAt:         200,
				EndAt:           400,
				RoomReservation: &RoomReservation{RoomID: testIDs[i], ReservationStatus: domain.ReservationStatus_Confirmed},
				ParticipantsIDs: []primitive.ObjectID{primitive.NewObjectID()},
				Notes:           &testStr,
				RemindAt:        0,
				IsDelete:        false,
			}
		} else {
			schedules = Event{
				ID:              primitive.NewObjectID(),
				Title:           testStr,
				Description:     &testStr,
				StartAt:         600,
				EndAt:           800,
				RoomReservation: &RoomReservation{RoomID: testIDs[i], ReservationStatus: domain.ReservationStatus_Confirmed},
				ParticipantsIDs: []primitive.ObjectID{primitive.NewObjectID()},
				Notes:           &testStr,
				RemindAt:        0,
				IsDelete:        false,
			}

		}

		_, err1 := testRoomScheduleRepo.roomScheduleCollection.InsertOne(context.TODO(), room)
		require.NoError(t, err1)
		_, err2 := testEventRepo.eventCollection.InsertOne(context.TODO(), schedules)
		require.NoError(t, err2)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	////////////////////////// Run test ////////////////////////////////////////////////////

	tests := []struct {
		name    string
		r       *mongoRoomScheduleRepository
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "Successful Upsert",
			r:    testRoomScheduleRepo, // Initialize with appropriate values
			args: args{
				ctx:        context.Background(), // Use context appropriate for testing
				roomIDs:    []string{testIDs[0].Hex(), testIDs[1].Hex(), testIDs[2].Hex()},
				equipments: []domain.Equipment{domain.EQUIPMENT_TABLE},
				rules:      []domain.Rule{domain.RULE_NO_FOOD},
				startAt:    0,
				endAt:      500,
				skip:       0,
				limit:      5,
			},
			want:    []int{3, 2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counts := []int{}
			eventCount := 0
			got, err := tt.r.QueryPaginated(tt.args.ctx, tt.args.roomIDs, tt.args.equipments, tt.args.rules, tt.args.startAt, tt.args.endAt, tt.args.skip, tt.args.limit)
			for _, roomSchedule := range got {
				eventCount += len(roomSchedule.Schedules)
			}
			counts = append(counts, len(got))
			counts = append(counts, eventCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoRoomScheduleRepository.QueryPaginated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(counts, tt.want) {
				t.Errorf("mongoRoomScheduleRepository.QueryPaginated() = %v, want %v", got, tt.want)
			}
		})
	}

	/////////////////// kill and remove the container //////////////////////////

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %s", err)
	}
}
