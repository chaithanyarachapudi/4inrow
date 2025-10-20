
package main

import (
    "fmt"

    "github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitKafkaProducer() error {
    brokers := []string{"localhost:29092"}
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    p, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return err
    }
    producer = p
    return nil
}

func EmitEvent(topic string, key string, value []byte) {
    if producer == nil {
        fmt.Println("kafka producer nil; skipping emit")
        return
    }
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(value),
    }
    _, _, err := producer.SendMessage(msg)
    if err != nil {
        fmt.Println("kafka send err:", err)
    }
}
