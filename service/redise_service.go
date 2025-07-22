package service

import (
	"context"
	"time"

	"test/pkg/logger"
	"test/storage"
)


type RedisService interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}


type redisService struct {
	redis storage.IRedisStorage
	log   logger.ILogger
}


func NewRedisService(stg storage.IStorage, log logger.ILogger, redis storage.IRedisStorage) RedisService {
	return &redisService{
		redis: redis,
		log:   log,
	}
}


func (s *redisService) SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	s.log.Info("RedisService.SetX called", logger.String("key", key))
	err := s.redis.SetX(ctx, key, value, duration)
	if err != nil {
		s.log.Error("Failed to set redis key", logger.Error(err))
	}
	return err
}


func (s *redisService) Get(ctx context.Context, key string) (string, error) {
	s.log.Info("RedisService.Get called", logger.String("key", key))
	val, err := s.redis.Get(ctx, key)
	if err != nil {
		s.log.Error("Failed to get redis key", logger.Error(err))
		return "", err
	}
	s.log.Info("Redis key fetched successfully", logger.String("key", key))
	return val, nil
}
