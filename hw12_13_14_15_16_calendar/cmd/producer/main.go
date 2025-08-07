package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/producer"
	sqlstorage "github.com/Calendar/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/producer.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()
	cfg, err := LoadProducerConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Load producer config: %v", *cfg)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	for attempt := 0; attempt < cfg.Producer.MaxRetries; attempt++ {
		err = producer.TryConnection(cfg.Broker)
		if err == nil {
			log.Printf("Successfully connected to broker\n")
			break
		}
		if attempt == cfg.Producer.MaxRetries-1 {
			log.Fatalf("Failed to connect to broker: %v", err)
		}
		backoff := float64(cfg.Producer.Interval) * math.Pow(2, float64(attempt))
		if backoff > cfg.Producer.MaxDelay {
			backoff = cfg.Producer.MaxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff))) * time.Second

		log.Printf("Attempt %d failed, waiting ~%v (max backoff: %v)...\n", attempt+1, jitter, backoff)
		time.Sleep(jitter)
	}

	kafkaProducer, err := producer.NewProducer(cfg.Broker, cfg.Producer)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()
	log.Printf("Successfully created Kafka producer\n")
	scanner := sqlstorage.NewScanner(cfg.DBConfig, ctx)
	err = scanner.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Successfully connected to database\n")
	ticker := time.NewTicker(time.Duration(cfg.Producer.CheckInterval) * time.Minute)
	defer ticker.Stop()
	tickerDelete := time.NewTicker(time.Duration(24) * time.Hour)
	defer tickerDelete.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Fatalf("The scanner is stopped by the context signal")
		case <-ticker.C:
			notification, err := scanner.GetNotificationTimeUntilEvent(ctx, cfg.Producer.CheckInterval)
			if err != nil {
				log.Printf("Failed to get event time until event: %v", err)
			}
			if notification == nil {
				log.Println("No notification time until event")
				continue
			}
			log.Printf("Notification time until event: %v", notification)
			for _, notice := range notification {
				jsonNotice, err := json.Marshal(notice)
				if err != nil {
					log.Printf("failed to serialize notification: %v", err)
					continue
				}
				log.Printf("Sending notification to Kafka producer: %v", string(jsonNotice))
				err = kafkaProducer.SendMessage(ctx, jsonNotice, cfg.Broker)
				if err != nil {
					log.Printf("Failed to write message to Kafka producer: %v", err)
				}
			}
			log.Printf("Successfully sent notification to Kafka producer\n")
		case <-tickerDelete.C:
			oneYear := time.Hour * 24 * 365
			err := scanner.DeleteNotificationTime(ctx, oneYear)
			if err != nil {
				log.Printf("event request error: %v", err)
			}
			log.Printf("Successfully deleted Kafka producer\n")
		}
	}

}
