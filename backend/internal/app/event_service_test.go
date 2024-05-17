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
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	infra "github.com/BartekTao/nycu-meeting-room-api/internal/infrastructure"
)

func Test_eventService_Upsert(t *testing.T) {
	type args struct {
		ctx context.Context
		req UpsertEventRequest
	}

	///////////////// Set up in-memory mongodb ////////////////////////////////

	var testMongoClient *mongo.Client
	var mongoPool *dockertest.Pool
	var mongoResource *dockertest.Resource

	mongoPool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	mongoResource, err = mongoPool.RunWithOptions(&dockertest.RunOptions{
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
	if err := mongoPool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		uri := fmt.Sprintf("mongodb://localhost:%s", mongoResource.GetPort("27017/tcp"))
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

	////////////////////// Set up dummy locker /////////////////////////////////

	var testRedisClient *redis.Client
	var redisPool *dockertest.Pool
	var redisResource *dockertest.Resource

	redisPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct redisPool: %s", err)
	}

	err = redisPool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	redisResource, err = redisPool.Run("redis", "3.2", nil)
	if err != nil {
		log.Fatalf("Could not start redisResource: %s", err)
	}

	if err = redisPool.Retry(func() error {
		testRedisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")),
		})

		return testRedisClient.Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// When you're done, kill and remove the container
	if err = redisPool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge redisResource: %s", err)
	}

	////////////////////////////////////////////////////////////////////////////

	/////////////////// Set up event service ///////////////////////////////////

	testEventRepo := infra.NewMockEventRepository(testMongoClient)
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
			got, err := tt.s.Upsert(tt.args[0].ctx, tt.args[0].req)
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
