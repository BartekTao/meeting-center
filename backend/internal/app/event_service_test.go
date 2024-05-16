package app

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	infra "github.com/BartekTao/nycu-meeting-room-api/internal/infrastructure"
)

var testMongoClient *mongo.Client
var pool *dockertest.Pool
var resource *dockertest.Resource

// SetupTestMongoDB starts an in-memory MongoDB container for testing.
func SetupTestMongoDB(t *testing.T) *mongo.Client {
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

	return testMongoClient
}

// TeardownTestMongoDB stops and removes the in-memory MongoDB container.
func TeardownTestMongoDB(t *testing.T) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %s", err)
	}
}

func Test_eventService_Upsert(t *testing.T) {
	type args struct {
		ctx context.Context
		req UpsertEventRequest
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

	/////////////////// Set up event service ///////////////////////////////////

	testEventRepo := infra.NewMongoEventRepository(testMongoClient)
	service := NewEventService(testEventRepo, locker)
	testEventService, ok := service.(*eventService)
	if !ok {
		log.Printf("Failed to type assert service to eventService")
		t.Skip("Skipping test due to type assertion failure")
		return
	}

	testStr := "Test"
	testID := "12341234123412341234AAAA"

	insertEvent := UpsertEventRequest{
		ID:              nil,
		Title:           testStr,
		Description:     &testStr,
		StartAt:         0,
		EndAt:           100,
		RoomID:          &testID,
		ParticipantsIDs: []string{testID},
		Notes:           &testStr,
		RemindAt:        0,
		UpdaterID:       testID,
	}

	////////////////////////////////////////////////////////////////////////////
	tests := []struct {
		name    string
		s       *eventService
		args    []args
		want    int
		wantErr bool
	}{
		{
			name: "Simultaneous upsert with same request",
			s:    testEventService, // Initialize with appropriate values
			args: []args{
				{
					ctx: context.Background(),
					req: insertEvent,
				},
				{
					ctx: context.Background(),
					req: insertEvent,
				},
				{
					ctx: context.Background(),
					req: insertEvent,
				},
				// Add more test cases as needed
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Upsert(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventService.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventService.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}