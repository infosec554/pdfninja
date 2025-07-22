package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type RoatateService interface {
	Create(ctx context.Context, req models.RotatePDFRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.RotateJob, error)
}

type rotateService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewRotateService(stg storage.IStorage, log logger.ILogger) RoatateService {
	return &rotateService{
		stg: stg,
		log: log,
	}
}

func (s *rotateService) Create(ctx context.Context, req models.RotatePDFRequest, userID string) (string, error) {
	s.log.Info("RotateService.Create called")

	// 1. Kiruvchi faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	jobID := uuid.NewString()
	outputPath := filepath.Join("storage/rotate", jobID+".pdf")

	job := &models.RotateJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		Angle:       req.Angle,
		Pages:       req.Pages,
		OutputPath:  outputPath,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.Rotate().Create(ctx, job); err != nil {
		s.log.Error("failed to create rotate job", logger.Error(err))
		return "", err
	}

	// 3. CLI orqali aylantirish
	// Masalan: pdfcpu rotate input.pdf output.pdf 90 1-3
	cmd := exec.Command("pdfcpu", "rotate", file.FilePath, outputPath, fmt.Sprintf("%d", req.Angle), req.Pages)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.log.Error("pdfcpu rotate failed", logger.Error(err))
		return "", err
	}

	// 4. Holatni yangilash
	job.Status = "done"
	if err := s.stg.Rotate().Update(ctx, job); err != nil {
		s.log.Error("failed to update rotate job", logger.Error(err))
		return "", err
	}

	s.log.Info("rotate job completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *rotateService) GetByID(ctx context.Context, id string) (*models.RotateJob, error) {
	job, err := s.stg.Rotate().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get rotate job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
