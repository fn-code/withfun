package main

import (
	"context"
	"github.com/fn-code/withfun/kafka"
	kc "github.com/segmentio/kafka-go"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	kcfg := &kafka.Config{
		Brokers:    []string{"localhost:9093"},
		GroupID:    "bank_account_microservice_consumer",
		InitTopics: true,
	}

	initKafkaTopics(ctx, kcfg)

	appSubs := NewBankAccountAppSubscription()
	consumerGroup := kafka.NewConsumerGroup([]string{"localhost:9093"}, "mongoGroup")

	err := consumerGroup.ConsumeTopicWithErrGroup(
		ctx,
		[]string{"eventstore_BankAccount"},
		10,
		appSubs.ProcessMessagesErrGroup,
	)
	if err != nil {
		log.Printf("(mongoConsumerGroup ConsumeTopicWithErrGroup) err: %v", err)
		cancel()
		return
	}

	select {}
}

func initKafkaTopics(ctx context.Context, kcfg *kafka.Config) {
	kafkaConn, err := kafka.NewKafkaConn(ctx, kcfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("(kafka connected) brokers: %+v\n", brokers)

	controller, err := kafkaConn.Controller()
	if err != nil {
		log.Fatal("kafka controller: ", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	log.Printf("(kafka controller uri) controllerURI: %s\n", controllerURI)

	conn, err := kc.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		log.Fatalf("initKafkaTopics.DialContext err: %v", err)
		os.Exit(1)
	}
	defer conn.Close() // nolint: errcheck

	log.Printf("(established new kafka controller connection) controllerURI: %s\n", controllerURI)

	bankAccountAggregateTopic := kc.TopicConfig{
		Topic:             "eventstore_BankAccount",
		NumPartitions:     10,
		ReplicationFactor: 1,
	}

	if err := conn.CreateTopics(bankAccountAggregateTopic); err != nil {
		log.Println("kafkaConn.CreateTopics", err)
		return
	}

	if err := conn.CreateTopics(bankAccountAggregateTopic); err != nil {
		log.Printf("(kafkaConn.CreateTopics) err: %v", err)
		return
	}

	log.Printf("(kafka topics created or already exists): %+v", []kc.TopicConfig{bankAccountAggregateTopic})
}
