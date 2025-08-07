package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/consumer"
	sqlstorage "github.com/Calendar/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/consumer.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()
	cfg, err := LoadConsumerConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Println(cfg.Broker.Topic)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	kafkaConsumer, err := consumer.NewConsumer(cfg.Broker)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()
	scanner := sqlstorage.NewScanner(cfg.DBConf, ctx)
	err = scanner.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")
	for {
		select {
		case <-ctx.Done():
			log.Println("The scanner is stopped by the context signal")
			return
		default:
			msg, err := kafkaConsumer.ReadMessage(ctx)
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				continue
			}
			err = scanner.CreateNotification(ctx, msg)
			if err != nil {
				log.Printf("Failed to create notification: %v", err)
			}
		}
	}
}
