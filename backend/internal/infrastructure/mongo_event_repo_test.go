package infra

import (
	"context"
	"reflect"
	"testing"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
)

func Test_mongoEventRepository_Upsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		event domain.Event
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
			m:    &mongoEventRepository{}, // Initialize with appropriate values
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
