package service

import (
	"context"
	"service-delivery/internal/cache"
	"service-delivery/internal/domain"
	"service-delivery/internal/repository/postgres"
)

type OrderService struct {
	repo  *postgres.OrderRepository
	cache *cache.OrderCache
}

func NewOrderService(repo *postgres.OrderRepository, cache *cache.OrderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *domain.Order) error {
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return err
	}
	s.cache.Set(order)
	return nil
}

func (s *OrderService) GetOrderByUID(ctx context.Context, orderUID string) (*domain.Order, error) {
	// Try to get from cache first
	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	// If not in cache, get from DB
	order, err := s.repo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	// If found in DB, add to cache
	if order != nil {
		s.cache.Set(order)
	}

	return order, nil
}

func (s *OrderService) RestoreCache(ctx context.Context) error {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}
	s.cache.Restore(orders)
	return nil
}
