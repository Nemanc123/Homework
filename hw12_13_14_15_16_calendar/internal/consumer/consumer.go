package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	produser_kafka "github.com/Calendar/hw12_13_14_15_calendar/internal/kafka"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

type Consumer struct {
	Reader *kafka.Reader
	Config *produser_kafka.BrokerConfig
}

func NewConsumer(cfg *produser_kafka.BrokerConfig) (*Consumer, error) {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{cfg.Broker},
			Topic:       cfg.Topic,
			GroupID:     "hw12_13_14_15_16_calendar",
			StartOffset: kafka.FirstOffset,
			Logger:      kafka.Logger(log.New(os.Stdout, "[KAFKA] ", log.LstdFlags)),
		}),
		Config: cfg,
	}, nil
}

func (p *Consumer) ReadMessage(ctx context.Context) (*storage.Notification, error) {
	var notification storage.Notification
	msg, err := p.Reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return &notification, err
	}
	fmt.Printf("msg: %+v\n", msg)
	err = json.Unmarshal(msg.Value, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (p *Consumer) Close() error {
	return p.Reader.Close()
}
