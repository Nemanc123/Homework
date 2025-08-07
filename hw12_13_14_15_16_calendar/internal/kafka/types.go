package kafka

type ProducerConfig struct {
	Interval      int     `yaml:"interval"`
	MaxRetries    int     `yaml:"maxRetries"`
	MaxDelay      float64 `yaml:"maxDelay"`
	CheckInterval int     `yaml:"checkInterval"`
	CheckDelete   int     `yaml:"checkDelete"`
}
type BrokerConfig struct {
	Broker string `yaml:"broker"`
	Topic  string `yaml:"topic"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	Database string `yaml:"database_db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
