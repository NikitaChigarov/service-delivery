package app

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"os/signal"
	"service-delivery/internal/cache"
	"service-delivery/internal/config"
	myhttp "service-delivery/internal/delivery/http"
	"service-delivery/internal/kafka"
	"service-delivery/internal/repository/postgres"
	"service-delivery/internal/service"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	httpServer *stdhttp.Server
	consumer   *kafka.Consumer
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize repository
	repo, err := postgres.NewOrderRepository(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	// Initialize cache
	orderCache := cache.NewOrderCache()

	// Initialize service
	orderService := service.NewOrderService(repo, orderCache)

	// Restore cache from database
	if err := orderService.RestoreCache(context.Background()); err != nil {
		return nil, err
	}

	// Initialize Kafka consumer
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, orderService)
	consumer.Start(context.Background())

	// Initialize HTTP server
	router := mux.NewRouter()
	myhttp.SetupRoutes(router, orderService)

	httpServer := &stdhttp.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	return &App{
		httpServer: httpServer,
		consumer:   consumer,
	}, nil
}

func (a *App) Run() error {
	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	// Close Kafka consumer
	if err := a.consumer.Close(); err != nil {
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
