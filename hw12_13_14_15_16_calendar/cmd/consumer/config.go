package main

import (
	"github.com/Calendar/hw12_13_14_15_calendar/internal/kafka"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Broker *kafka.BrokerConfig `yaml:"kafka"`
	DBConf *kafka.DBConfig     `yaml:"database"`
}

func LoadConsumerConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
