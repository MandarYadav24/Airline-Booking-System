package kafka

import (
	"log"
	"time"

	"airline-booking/pkg/config"

	"github.com/IBM/sarama"
)

type Producer struct {
	Client sarama.SyncProducer
}

// NewProducer initializes Kafka producer
func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll
	kafkaCfg.Producer.Retry.Max = cfg.Producer.Retries
	kafkaCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaCfg.Net.DialTimeout = 10 * time.Second
	kafkaCfg.Version = sarama.V3_4_0_0 // compatible with Kafka 4.x (KRaft)

	producer, err := sarama.NewSyncProducer(cfg.Brokers, kafkaCfg)
	if err != nil {
		return nil, err
	}
	//log.Println("Kafka producer connected to:", cfg.Brokers)
	return &Producer{Client: producer}, nil
}

// SendMessage publishes a message to a given topic
func (p *Producer) SendMessage(topic string, key string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.Client.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to topic %s [partition=%d, offset=%d]", topic, partition, offset)
	return nil
}

func (p *Producer) Close() {
	if p.Client != nil {
		_ = p.Client.Close()
		log.Println("Kafka producer closed")
	}
}
