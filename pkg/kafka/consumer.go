package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"

	"airline-booking/pkg/config"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Group sarama.ConsumerGroup
}

func NewConsumer(cfg *config.KafkaConfig) (*Consumer, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	kafkaCfg.Version = sarama.V3_4_0_0 // compatible with Kafka 4.x (KRaft)

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, kafkaCfg)
	if err != nil {
		return nil, err
	}

	log.Println("Kafka consumer connected with group:", cfg.GroupID)
	return &Consumer{Group: group}, nil
}

// RunConsumer listens to Kafka messages
func (c *Consumer) RunConsumer(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	go func() {
		for {
			if err := c.Group.Consume(ctx, topics, handler); err != nil {
				log.Printf(" Error consuming: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm

	log.Println("Stopping Kafka consumer...")
	c.Close()
}

func (c *Consumer) Close() {
	if c.Group != nil {
		if err := c.Group.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		} else {
			log.Println("Kafka consumer closed cleanly")
		}
	}
}

type ExampleHandler struct{}

func (ExampleHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer setup complete")
	return nil
}
func (ExampleHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer cleanup done")
	return nil
}
func (ExampleHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message: topic=%s partition=%d offset=%d value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
