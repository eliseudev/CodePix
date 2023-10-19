package kafka

import (
	"codepix/application/factory"
	appmodel "codepix/application/model"
	"codepix/application/usecase"
	"codepix/domain/model"
	"fmt"
	fkafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
	"os"
)

type KafkaProcessor struct {
	Database     *gorm.DB
	Producer     *fkafka.Producer
	DeliveryChan chan fkafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *fkafka.Producer, deliveryChan chan fkafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &fkafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := fkafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	_ = c.SubscribeTopics(topics, nil)

	fmt.Println("Kafka consumer has been started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *fkafka.Message) {
	transactionTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionTopic:
		_ = k.processTransaction(msg)
	case transactionConfirmationTopic:
		_ = k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *fkafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description)
	if err != nil {
		fmt.Println("error register transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *fkafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)
	if transaction.Status == model.TransactionConfirmed {
		err := k.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err == nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}
