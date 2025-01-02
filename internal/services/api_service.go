package services

import (
	"github.com/ronaldalds/base-go-api/internal/database"
)

type HealthHandler struct {
	Sql   map[string]string `json:"sql"`
	Redis map[string]string `json:"redis"`
}

type Service struct {
	GormStore  *database.GormStore
	RedisStore *database.RedisStore
}

func NewService() *Service {
	return &Service{
		GormStore:  database.DB.GormStore,
		RedisStore: database.DB.RedisStore,
	}
}

func (s *Service) Health() *HealthHandler {
	stats := &HealthHandler{
		Sql:   s.GormStore.Health(),
		Redis: s.RedisStore.Health(),
	}
	return stats
}
