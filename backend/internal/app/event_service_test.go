package app

import (
	"context"
	"errors"
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

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	infra "github.com/BartekTao/nycu-meeting-room-api/internal/infrastructure"
)

type MockEventRepository struct {
	infra.BaseRepository[infra.Event]
	client          *mongo.Client
	eventCollection *mongo.Collection
}

func NewTestEventRepository(client *mongo.Client) domain.EventRepository {
	return MockEventRepository{
		client:          client,
		eventCollection: client.Database("meetingCenter").Collection("events"),
	}
}

func (m *MockEventRepository) Upsert(ctx context.Context, event domain.Event) (*domain.Event, error) {
	if event.ID == "" {
		return nil, errors.New("event ID cannot be empty")
	}
	// Store or update the event in the repository
	m.Events[event.ID] = event
	return &event, nil
}

// Implement the Delete method of the EventRepository interface
func (m *MockEventRepository) Delete(ctx context.Context, id string) (*domain.Event, error) {
	dummy := domain.Event{}
	return &dummy, nil
}

// Implement the UpdateSummary method of the EventRepository interface
func (m *MockEventRepository) UpdateSummary(ctx context.Context, id string, summary string, updaterID string) (bool, error) {
	return true, nil
}

// Implement the GetByID method of the EventRepository interface
func (m *MockEventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	dummy := domain.Event{}
	return &dummy, nil
}

// Implement the GetByUsers method of the EventRepository interface
func (m *MockEventRepository) GetByUsers(ctx context.Context, ids []string, startAt, endAt int64) (map[string][]domain.Event, error) {
	dummyUserEvents := make(map[string][]domain.Event)
	return dummyUserEvents, nil
}

// Implement the CheckAvailableRoom method of the EventRepository interface
func (m *MockEventRepository) CheckAvailableRoom(ctx context.Context, roomID string, startAt, endAt int64) (bool, error) {
	// Implement logic to check room availability for a given time range
	// For testing purposes, return a mock result (always true for now)
	return true, nil
}

// Implement the GetAllWithRoomConfirmed method of the EventRepository interface
func (m *MockEventRepository) GetAllWithRoomConfirmed(ctx context.Context, roomIDs []string, startAt, endAt int64) ([]domain.Event, error) {
	var dummyEvents []domain.Event
	return dummyEvents, nil
}

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

	redisPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct redisPool: %s", err)
	}

	err = redisPool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	redisResource, err := redisPool.Run("redis", "3.2", nil)
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
