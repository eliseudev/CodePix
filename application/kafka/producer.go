package kafka

import (
	"fmt"
	fkafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func NewKafkaProducer() *fkafka.Producer {
	confMap := &fkafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
	}
	p, err := fkafka.NewProducer(confMap)
	if err != nil {
		panic(err)
	}

	return p
}

func Publish(msg string, topic string, producer *fkafka.Producer, deliveryChanel chan fkafka.Event) error {
	message := &fkafka.Message{
		TopicPartition: fkafka.TopicPartition{Topic: &topic, Partition: fkafka.PartitionAny},
		Value:          []byte(msg),
	}
	err := producer.Produce(message, deliveryChanel)
	if err != nil {
		return err
	}

	return nil
}

func DeliveryReport(deliveryChanel chan fkafka.Event) {
	for e := range deliveryChanel {
		switch ev := e.(type) {
		case *fkafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery failed:", ev.TopicPartition)
			} else {
				fmt.Println("Delivery message to:", ev.TopicPartition)
			}
		}
	}
}
