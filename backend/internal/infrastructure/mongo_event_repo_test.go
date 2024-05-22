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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_mongoEventRepository_Upsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		event domain.Event
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

	////////////////////////////////////////////////////////////////////////////

	/////////////////// Set up event repo //////////////////////////////////////

	repo := NewMongoEventRepository(testMongoClient)
	testEventRepo, ok := repo.(*mongoEventRepository)
	if !ok {
		log.Printf("Failed to type assert repo to mongoEventRepository")
		t.Skip("Skipping test due to type assertion failure")
		return
	}

	testStr := "Test"
	testID := "12341234123412341234AAAA"

	insert_event := domain.Event{
		ID:              nil,
		Title:           testStr,
		Description:     &testStr,
		StartAt:         0,
		EndAt:           0,
		ParticipantsIDs: []string{testID},
		Notes:           &testStr,
		RemindAt:        0,
		UpdaterID:       testID,
	}

	var roomReservation *domain.RoomReservation = &domain.RoomReservation{
		RoomID:            &testID,
		ReservationStatus: domain.ReservationStatus_Confirmed,
	}
	insert_event.RoomReservation = roomReservation

	////////////////////////////////////////////////////////////////////////////

	/////////////////////////// Run tests //////////////////////////////////////

	tests := []struct {
		name    string
		m       *mongoEventRepository
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Successful Upsert",
			m:    testEventRepo, // Initialize with appropriate values
			args: args{
				ctx:   context.Background(), // Use context appropriate for testing
				event: insert_event,
			},
			want:    domain.ReservationStatus_Confirmed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upsert(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoEventRepository.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.RoomReservation.ReservationStatus, tt.want) {
				t.Errorf("mongoEventRepository.Upsert() = %v, want %v", got.RoomReservation, tt.want)
			}
		})
	}

	////////////////////////////////////////////////////////////////////////////

	/////////////////// kill and remove the container //////////////////////////

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %s", err)
	}
}
