/*
Copyright © 2023 343.IT <eliseu.dev@outlook.com>
*/
package cmd

import (
	"codepix/application/kafka"
	"codepix/infrastructure/db"
	fkafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start consuming transaction using Apache Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		deliveryChan := make(chan fkafka.Event)
		database := db.ConnectDB("env")
		producer := kafka.NewKafkaProducer()

		go kafka.DeliveryReport(deliveryChan)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
