package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	//log     logger.Logger
	brokers []string
	w       *kafka.Writer
}

func Errorf(s string, args ...interface{}) {
	log.Printf(s, args...)
}

// NewProducer create new kafka producer
func NewProducer(brokers []string) *producer {
	return &producer{brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(Errorf))}
}

// NewAsyncProducer create new kafka producer
func NewAsyncProducer(brokers []string) *producer {
	return &producer{brokers: brokers, w: NewAsyncWriter(brokers, kafka.LoggerFunc(Errorf))}
}

// NewAsyncProducerWithCallback create new kafka producer with callback for delete invalid projection
func NewAsyncProducerWithCallback(brokers []string, cb AsyncWriterCallback) *producer {
	return &producer{brokers: brokers, w: NewAsyncWriterWithCallback(brokers, kafka.LoggerFunc(Errorf), cb)}
}

// NewRequireNoneProducer create new fire and forget kafka producer
func NewRequireNoneProducer(brokers []string) *producer {
	return &producer{brokers: brokers, w: NewRequireNoneWriter(brokers, kafka.LoggerFunc(Errorf))}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "producer.PublishMessage")
	//defer span.Finish()

	if err := p.w.WriteMessages(ctx, msgs...); err != nil {
		//tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.w.Close()
}
