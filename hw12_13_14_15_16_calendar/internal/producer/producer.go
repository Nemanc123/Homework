package producer

import (
	"context"
	produser_kafka "github.com/Calendar/hw12_13_14_15_calendar/internal/kafka"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	Writer       *kafka.Writer
	Config       *produser_kafka.ProducerConfig
	BrokerConfig *produser_kafka.BrokerConfig
}

func TryConnection(cfg *produser_kafka.BrokerConfig) error {
	_, err := kafka.Dial("tcp", cfg.Broker)
	if err != nil {
		return err
	}
	return nil
}
func NewProducer(cfg *produser_kafka.BrokerConfig, brk *produser_kafka.ProducerConfig) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:        kafka.TCP(cfg.Broker),
		Topic:       cfg.Topic,
		MaxAttempts: 5,
	}
	return &Producer{
		Writer:       writer,
		Config:       brk,
		BrokerConfig: cfg,
	}, nil
}
func (p *Producer) SendMessage(ctx context.Context, jsonNotice []byte, cfg *produser_kafka.BrokerConfig) error {

	err := p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: jsonNotice,
	})
	if err != nil {
		return err
	}
	return nil
}
func (p *Producer) Close() error {
	return p.Writer.Close()
}
