package kafka

import (
	"context"
	"log"
	"route256/loms/internal/service"
	"time"

	"github.com/Shopify/sarama"
)

type Sender interface {
	SendNotification(ctx context.Context, msg service.OutboxMessage) error
}

type sender struct {
	topic    string
	producer sarama.SyncProducer
}

func NewSender(brokers []string, topic string) (Sender, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, err
	}
	topics, err := admin.ListTopics()
	if err != nil {
		return nil, err
	}
	if _, ok := topics[topic]; !ok {
		err = admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     3,
			ReplicationFactor: 3,
		}, false)
		if err != nil {
			return nil, err
		}
	}

	config.Producer.Partitioner = sarama.NewHashPartitioner // sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true // exactly once
	config.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &sender{
		topic:    topic,
		producer: producer,
	}, nil
}

func (s sender) SendNotification(ctx context.Context, msg service.OutboxMessage) error {
	m := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: -1,
		Value:     sarama.StringEncoder(msg.Message),
		Key:       sarama.StringEncoder(msg.Key),
		Timestamp: time.Now(),
	}

	partition, offset, err := s.producer.SendMessage(m)
	if err != nil {
		return err
	}

	log.Printf("notification sent for order %v, message %v, partition %v, offset %v", msg.Key, msg.Message, partition, offset)
	return nil
}
