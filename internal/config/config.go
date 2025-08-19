package config

import "os"

type Config struct {
	HTTPAddr      string
	PostgresDSN   string
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaGroupID  string
}

func Load() (*Config, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://orders_user:orders_password@localhost:5432/orders_db?sslmode=disable"
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "orders"
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "order-service-group"
	}

	return &Config{
		HTTPAddr:      httpAddr,
		PostgresDSN:   dsn,
		KafkaBrokers:  []string{kafkaBrokers},
		KafkaTopic:    kafkaTopic,
		KafkaGroupID:  kafkaGroupID,
	}, nil
}