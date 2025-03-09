package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type appSubscription struct {
}

func NewBankAccountAppSubscription() *appSubscription {
	return &appSubscription{}
}

func (s *appSubscription) ProcessMessagesErrGroup(ctx context.Context, r *kafka.Reader, workerID int) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Printf("(mongoSubscription) workerID: %d, err: %v\n", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case "eventstore_BankAccount":
			s.handleBankAccountEvents(ctx, r, m)
		}
	}
}

// EventType is the type of any event, used as its unique identifier.
type EventType string

// AggregateType type of the Aggregate
type AggregateType string

// Event is an internal representation of an event, returned when the Aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type Event struct {
	EventID       string
	AggregateID   string
	EventType     EventType
	AggregateType AggregateType
	Version       uint64
	Data          []byte
	Metadata      []byte
	Timestamp     time.Time
}

func (s *appSubscription) handleBankAccountEvents(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	log.Println("appSubscription.handleBankAccountEvents")

	var events []Event
	if err := json.Unmarshal(m.Value, &events); err != nil {
		log.Printf("serializer.Unmarshal: %v", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	for _, event := range events {
		//if err := s.handle(ctx, r, m, event); err != nil {
		//	return
		//}
		log.Println("-> Receive event : ", event)
	}
	s.commitMessage(ctx, r, m)
}

func (s *appSubscription) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	if err := r.CommitMessages(ctx, m); err != nil {
		log.Printf("(mongoSubscription) [CommitMessages] err: %v", err)
		return
	}
	log.Println("(Committed Kafka message)", m.Topic, m.Partition, m.Offset)
}

func (s *appSubscription) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	if err := r.CommitMessages(ctx, m); err != nil {
		log.Printf("(mongoSubscription) [CommitMessages] err: %v", err)
		return
	}
	log.Println("(Committed Kafka message)", m.Topic, m.Partition, m.Offset)
}

func (s *appSubscription) logProcessMessage(m kafka.Message, workerID int) {
	log.Println("(Processing Kafka message)", m.Topic, m.Partition, m.Value, workerID, m.Offset, m.Time)
}
