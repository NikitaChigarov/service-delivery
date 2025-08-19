package cache

import (
	"service-delivery/internal/domain"
	"sync"
)

type OrderCache struct {
	mu    sync.RWMutex
	cache map[string]*domain.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		cache: make(map[string]*domain.Order),
	}
}

func (c *OrderCache) Set(order *domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[order.OrderUID] = order
}

func (c *OrderCache) Get(orderUID string) (*domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.cache[orderUID]
	return order, exists
}

func (c *OrderCache) Restore(orders []domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		orderCopy := order
		c.cache[order.OrderUID] = &orderCopy
	}
}
