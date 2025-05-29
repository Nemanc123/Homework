package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	sqlstorage "github.com/Calendar/hw12_13_14_15_calendar/internal/storage/sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf              `yaml:"logger"`
	Database sqlstorage.DatabaseConf `yaml:"database"`
	Storage  string                  `yaml:"storage_impl"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

func NewConfig() Config {
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return *config
}
func loadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}
	var ConfigFile Config
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	err = yaml.Unmarshal(file, &ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &ConfigFile, nil
}
