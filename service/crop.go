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

type CropPDFService interface {
	Create(ctx context.Context, req models.CropPDFRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.CropPDFJob, error)
}

type cropPDFService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewCropPDFService(stg storage.IStorage, log logger.ILogger) CropPDFService {
	return &cropPDFService{stg: stg, log: log}
}

func (s *cropPDFService) Create(ctx context.Context, req models.CropPDFRequest, userID string) (string, error) {
	s.log.Info("CropPDFService.Create called")

	// 1. Kirish faylini olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", fmt.Errorf("file not found")
	}

	// 2. Job yaratish
	jobID := uuid.New().String()
	outputFileID := uuid.New().String()
	outputPath := filepath.Join("storage/crop_pdf", outputFileID+".pdf")

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		s.log.Error("failed to create output directory", logger.Error(err))
		return "", err
	}

	job := &models.CropPDFJob{
		ID:           jobID,
		UserID:       userID,
		InputFileID:  req.InputFileID,
		OutputFileID: outputFileID,
		Box:          req.Box,
		Pages:        req.Pages,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.stg.Crop().Create(ctx, job); err != nil {
		s.log.Error("failed to create crop job", logger.Error(err))
		return "", err
	}

	// 3. PDFCPU orqali crop qilish
	args := []string{
		"pages",
		"trim",
		"-pages", req.Pages,
		"-mediabox", req.Box,
		file.FilePath,
		outputPath,
	}
	cmd := exec.Command("pdfcpu", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.log.Error("pdfcpu crop failed", logger.Error(err))
		return "", err
	}

	// 4. Yangi faylni saqlash
	fi, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("cannot stat output file", logger.Error(err))
		return "", err
	}

	newFile := models.File{
		ID:         outputFileID,
		UserID:     userID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   fi.Size(),
		UploadedAt: time.Now(),
	}

	if _, err := s.stg.File().Save(ctx, newFile); err != nil {
		s.log.Error("failed to save output file", logger.Error(err))
		return "", err
	}

	// 5. Holatni yangilash
	job.Status = "done"
	if err := s.stg.Crop().Update(ctx, job); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("Crop PDF job completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *cropPDFService) GetByID(ctx context.Context, id string) (*models.CropPDFJob, error) {
	job, err := s.stg.Crop().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get crop job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
