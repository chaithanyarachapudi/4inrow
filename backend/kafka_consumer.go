package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/IBM/sarama"
)

// StartKafkaConsumer starts the Kafka consumer as a goroutine
func StartKafkaConsumer(brokers []string, group string, topic string) {
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

    client, err := sarama.NewConsumerGroup(brokers, group, config)
    if err != nil {
        log.Fatal("Failed to create consumer group:", err)
    }

    go func() {
        defer client.Close()

        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        // Handle OS signals to gracefully shut down
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
        go func() {
            <-sig
            cancel()
        }()

        handler := ConsumerGroupHandler{}
        for {
            if err := client.Consume(ctx, []string{topic}, handler); err != nil {
                fmt.Println("Consume error:", err)
            }
            if ctx.Err() != nil {
                fmt.Println("Kafka consumer shutting down...")
                return
            }
        }
    }()
}

// ConsumerGroupHandler handles Kafka messages
type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        fmt.Printf("Analytics received: topic=%s key=%s value=%s\n",
            msg.Topic, string(msg.Key), string(msg.Value))
        sess.MarkMessage(msg, "")
    }
    return nil
}
