package main

import (
	"context"
	"fmt"
	"log"
	"time"

	infra "github.com/BartekTao/nycu-meeting-room-api/internal/infrastructure"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/notification"
)

func main() {
	ctx := context.Background()

	now := time.Now().UTC()
	checkAt := getNearest5MinUnixMs(now)

	mongoClient := infra.SetUpMongoDB()
	defer infra.ShutdownMongoDB(mongoClient)

	eventRepo := infra.NewMongoEventRepository(mongoClient)
	userRepo := infra.NewMongoUserRepo(mongoClient)

	mailHandler, err := notification.NewGmailHandler()
	if err != nil {
		log.Fatalf("Error initializing mail handler: %v", err)
	}

	events, err := eventRepo.GetRemindEvents(ctx, checkAt)
	if err != nil {
		log.Fatalf("Error getting events: %v", err)
	}

	userIDs := make([]string, 0)
	for _, event := range events {
		userIDs = append(userIDs, event.ParticipantsIDs...)
	}

	users, err := userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}

	userEmailMap := make(map[string]string)
	for _, user := range users {
		userEmailMap[*user.ID] = user.Email
	}

	for _, event := range events {
		var recipientEmails []string
		for _, userID := range event.ParticipantsIDs {
			email, exists := userEmailMap[userID]
			if exists {
				recipientEmails = append(recipientEmails, email)
			}
		}

		startTime := time.UnixMilli(event.StartAt)
		loc, err := time.LoadLocation("Local")
		if err != nil {
			log.Printf("Error loading location: %v\n", err)
			continue
		}
		startTime = startTime.In(loc)
		formattedTime := startTime.Format("2006-01-02 15:04:05 -0700")
		subject := fmt.Sprintf("Reminder: %s - %s", event.Title, formattedTime)
		content := fmt.Sprintf("The meeting will start at %s", formattedTime)

		err = mailHandler.Send(recipientEmails, subject, content)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			continue
		}
	}
}

func getNearest5MinUnixMs(t time.Time) int64 {
	minutes := t.Minute() % 5
	seconds := t.Second()
	nanoseconds := t.Nanosecond()

	roundedTime := t.Add(-time.Duration(minutes) * time.Minute).
		Add(-time.Duration(seconds) * time.Second).
		Add(-time.Duration(nanoseconds) * time.Nanosecond)

	return roundedTime.UnixMilli()
}
