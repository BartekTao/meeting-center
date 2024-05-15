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

// Example usage in test functions:
func TestMyMongoDBFunction(t *testing.T) {
	client := SetupTestMongoDB(t)
	defer TeardownTestMongoDB(t)

	// Use the testMongoClient to interact with MongoDB in your test function
	// For example:
	ctx := context.Background()
	collection := client.Database("testdb").Collection("mycollection")

	// Perform test operations on MongoDB collection
	_, err := collection.InsertOne(ctx, bson.M{"name": "test document"})
	if err != nil {
		t.Fatalf("Failed to insert document into MongoDB: %s", err)
	}

	// Add more test assertions as needed
}

func Test_mongoEventRepository_Upsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		event domain.Event
	}

	client := SetupTestMongoDB(t)
	defer TeardownTestMongoDB(t)

	repo := NewMongoEventRepository(client)
	testEventRepo, ok := repo.(*mongoEventRepository)
	if !ok {
		t.Fatal("Failed to type assert repo to mongoEventRepository")
	}

	tests := []struct {
		name    string
		m       *mongoEventRepository
		args    args
		want    *domain.Event
		wantErr bool
	}{
		{
			name: "Successful Upsert",
			m:    testEventRepo, // Initialize with appropriate values
			args: args{
				ctx:   context.Background(), // Use context appropriate for testing
				event: domain.Event{ /* Populate request with test data */ },
			},
			want:    &domain.Event{ /* Define expected result */ },
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoEventRepository.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoEventRepository_Upsert_multiple(t *testing.T) {
	type args struct {
		ctx   context.Context
		event domain.Event
	}
	tests := []struct {
		name    string
		m       *mongoEventRepository
		args    []args
		want    int
		wantErr bool
	}{
		{
			name: "Upsert multiple events with same room and time",
			m:    &mongoEventRepository{}, // Initialize with appropriate values
			args: []args{
				{
					ctx:   context.Background(),
					event: domain.Event{},
				},
				{
					ctx:   context.Background(),
					event: domain.Event{},
				},
				{
					ctx:   context.Background(),
					event: domain.Event{},
				},
				// Add more test cases as needed
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upsert(tt.args[0].ctx, tt.args[0].event)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoEventRepository.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoEventRepository.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}
