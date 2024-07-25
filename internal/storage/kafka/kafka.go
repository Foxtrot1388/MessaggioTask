package kafka

import (
	"encoding/json"
	"log/slog"

	"github.com/Foxtrot1388/MessaggioTask/internal/config"
	"github.com/IBM/sarama"
)

type kafkaStorage struct {
	log      *slog.Logger
	cfg      *config.AppConfig
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

type messageKafka struct {
	ID int
}

func New(cfg *config.AppConfig, log *slog.Logger) (*kafkaStorage, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokers := []string{cfg.KafkaHost}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	return &kafkaStorage{log: log, cfg: cfg, producer: producer, consumer: consumer}, nil

}

func (con *kafkaStorage) Send(ID int) error {

	message := messageKafka{
		ID: ID,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: con.cfg.KafkaTopic,
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = con.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil

}

func (con *kafkaStorage) Read() (<-chan int, error) {

	consumer, err := con.consumer.ConsumePartition(con.cfg.KafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	ch := make(chan int)
	go func() {
		defer close(ch)
		for msg := range consumer.Messages() {
			var notification messageKafka
			err := json.Unmarshal(msg.Value, &notification)
			if err != nil {
				con.log.Error(err.Error())
				continue
			}
			ch <- notification.ID
		}
	}()

	return ch, nil

}
