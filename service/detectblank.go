package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/detectblank"
	"test/pkg/logger"
	"test/storage"
)

type DetectBlankService interface {
	Create(ctx context.Context, inputFileID string, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.DetectBlankPagesJob, error)
}

type detectBlankService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewDetectBlankService(stg storage.IStorage, log logger.ILogger) DetectBlankService {
	return &detectBlankService{stg: stg, log: log}
}

func (s *detectBlankService) Create(ctx context.Context, inputFileID string, userID string) (string, error) {
	s.log.Info("DetectBlankService.Create called")

	file, err := s.stg.File().GetByID(ctx, inputFileID)
	if err != nil {
		return "", fmt.Errorf("input file not found")
	}

	jobID := uuid.New().String()
	job := &models.DetectBlankPagesJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: inputFileID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.DetectBlankPages().Create(ctx, job); err != nil {
		return "", err
	}

	blankPages, err := detectblank.DetectBlankPages(file.FilePath, 10) // masalan, 10 — minimal matn uzunligi

	if err != nil {
		s.log.Error("failed to detect blank pages", logger.Error(err))
		job.Status = "failed"
		_ = s.stg.DetectBlankPages().Update(ctx, job)
		return "", err
	}

	job.BlankPages = blankPages
	job.Status = "done"

	if err := s.stg.DetectBlankPages().Update(ctx, job); err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *detectBlankService) GetByID(ctx context.Context, id string) (*models.DetectBlankPagesJob, error) {
	return s.stg.DetectBlankPages().GetByID(ctx, id)
}
