package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type SharedLinkService interface {
	Create(ctx context.Context, req models.CreateSharedLinkRequest) (string, error)
	GetByToken(ctx context.Context, token string) (*models.SharedLink, error)
}

type sharedLinkService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewSharedLinkService(stg storage.IStorage, log logger.ILogger) SharedLinkService {
	return &sharedLinkService{stg: stg, log: log}
}

func (s *sharedLinkService) Create(ctx context.Context, req models.CreateSharedLinkRequest) (string, error) {
	s.log.Info("SharedLinkService.Create called")

	linkID := uuid.New().String()
	token := uuid.New().String() // token sifatida UUID ishlatyapmiz, kerak bo‘lsa short-link generator bilan almashtirish mumkin

	link := &models.SharedLink{
		ID:          linkID,
		FileID:      req.FileID,
		SharedToken: token,
		ExpiresAt:   req.ExpiresAt,
		CreatedAt:   time.Now(),
	}

	if err := s.stg.SharedLink().Create(ctx, link); err != nil {
		s.log.Error("failed to create shared link", logger.Error(err))
		return "", err
	}

	return token, nil
}

func (s *sharedLinkService) GetByToken(ctx context.Context, token string) (*models.SharedLink, error) {
	s.log.Info("SharedLinkService.GetByToken called")
	return s.stg.SharedLink().GetByToken(ctx, token)
}
