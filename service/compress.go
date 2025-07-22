package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type CompressService interface {
	Create(ctx context.Context, req models.CompressRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.CompressJob, error)
}

type compressService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewCompressService(stg storage.IStorage, log logger.ILogger) CompressService {
	return &compressService{
		stg: stg,
		log: log,
	}
}

func (s *compressService) Create(ctx context.Context, req models.CompressRequest, userID string) (string, error) {
	s.log.Info("CompressService.Create called")

	// 1. Kiruvchi faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	job := &models.CompressJob{
		ID:          uuid.NewString(),
		UserID:      userID,
		InputFileID: req.InputFileID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.Compress().Create(ctx, job); err != nil {
		s.log.Error("failed to create compress job", logger.Error(err))
		return "", err
	}

	// 3. PDF ni siqish (CLI orqali pdfcpu)
	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/compress", outputID+".pdf")

	cmd := exec.Command("pdfcpu", "optimize", file.FilePath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.log.Error("pdfcpu optimize failed", logger.Error(err))
		return "", err
	}

	// 4. Output faylni metadata bilan DB ga yozish
	info, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("cannot stat output file", logger.Error(err))
		return "", err
	}

	newFile := models.File{
		ID:         outputID,
		UserID:     userID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   info.Size(),
		UploadedAt: time.Now(),
	}

	if _, err := s.stg.File().Save(ctx, newFile); err != nil {
		s.log.Error("failed to save output file", logger.Error(err))
		return "", err
	}

	// 5. Job yangilash
	job.OutputFileID = outputID
	job.Status = "done"

	if err := s.stg.Compress().Update(ctx, job); err != nil {
		s.log.Error("failed to update compress job", logger.Error(err))
		return "", err
	}

	s.log.Info("compress job completed", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *compressService) GetByID(ctx context.Context, id string) (*models.CompressJob, error) {
	job, err := s.stg.Compress().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get compress job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
