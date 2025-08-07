package main

import (
	"github.com/Calendar/hw12_13_14_15_calendar/internal/kafka"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Broker   *kafka.BrokerConfig   `yaml:"kafka"`
	Producer *kafka.ProducerConfig `yaml:"producer"`
	DBConfig *kafka.DBConfig       `yaml:"database"`
}

func LoadProducerConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Устанавливаем значения по умолчанию
	if cfg.Producer.Interval == 0 {
		cfg.Producer.Interval = 1
	}

	return &cfg, nil
}
