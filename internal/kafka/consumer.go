package kafka

import (
	"context"
	"encoding/json"
	"log"
	"service-delivery/internal/domain"
	"service-delivery/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
}

func NewConsumer(brokers []string, topic string, groupID string, service *service.OrderService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					log.Printf("error fetching message: %v", err)
					continue
				}

				var order domain.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("error unmarshalling message: %v", err)
					continue
				}

				if err := c.service.ProcessOrder(ctx, &order); err != nil {
					log.Printf("error processing order: %v", err)
					continue
				}

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("error committing message: %v", err)
				}
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
