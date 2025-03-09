package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

// WorkerErrGroup kafka consumer worker fetch and process messages from reader
type WorkerErrGroup func(ctx context.Context, r *kafka.Reader, workerID int) error

type consumerGroup struct {
	Brokers []string
	GroupID string
	//log     logger.Logger
}

// NewConsumerGroup kafka consumer group constructor
func NewConsumerGroup(brokers []string, groupID string) *consumerGroup {
	return &consumerGroup{Brokers: brokers, GroupID: groupID}
}

// NewKafkaReader create new kafka reader
func (c *consumerGroup) NewKafkaReader(kafkaURL []string, groupTopics []string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		GroupTopics:            groupTopics,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		MaxAttempts:            maxAttempts,
		MaxWait:                maxWait,
		Dialer:                 &kafka.Dialer{Timeout: dialTimeout},
	})
}

// ConsumeTopic start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker) {
	r := c.NewKafkaReader(c.Brokers, groupTopics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			//c.log.Warnf("consumerGroup.r.Close: %v", err)
			log.Printf("consumerGroup.r.Close: %v", err)
		}
	}()

	//c.log.Infof("(Starting consumer groupID): GroupID %s, topic: %+v, poolSize: %v", c.GroupID, groupTopics, poolSize)
	log.Printf("(Starting consumer groupID): GroupID %s, topic: %+v, poolSize: %v", c.GroupID, groupTopics, poolSize)

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}

// ConsumeTopicWithErrGroup start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopicWithErrGroup(ctx context.Context, groupTopics []string, poolSize int, worker WorkerErrGroup) error {
	r := c.NewKafkaReader(c.Brokers, groupTopics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			//c.log.Warnf("consumerGroup.r.Close: %v", err)
			log.Printf("consumerGroup.r.Close: %v", err)
		}
	}()

	//c.log.Infof("(Starting ConsumeTopicWithErrGroup) GroupID: %s, topics: %+v, poolSize: %d", c.GroupID, groupTopics, poolSize)

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= poolSize; i++ {
		g.Go(c.runWorker(ctx, worker, r, i))
	}
	return g.Wait()
}

func (c *consumerGroup) runWorker(ctx context.Context, worker WorkerErrGroup, r *kafka.Reader, i int) func() error {
	return func() error {
		return worker(ctx, r, i)
	}
}
