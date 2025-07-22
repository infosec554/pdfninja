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

type ProtectPDFService interface {
	Create(ctx context.Context, req models.ProtectPDFRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.ProtectPDFJob, error)
}

type protectPDFService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewProtectPDFService(stg storage.IStorage, log logger.ILogger) ProtectPDFService {
	return &protectPDFService{
		stg: stg,
		log: log,
	}
}

func (s *protectPDFService) Create(ctx context.Context, req models.ProtectPDFRequest, userID string) (string, error) {
	s.log.Info("ProtectPDFService.Create called")

	// 1. Faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", fmt.Errorf("file not found")
	}

	// 2. Output joyini tayyorlash
	jobID := uuid.NewString()
	outputFileID := uuid.NewString()
	outputPath := filepath.Join("storage/protect_pdf", outputFileID+".pdf")

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		s.log.Error("failed to create output directory", logger.Error(err))
		return "", err
	}

	// 3. Job yaratish
	job := &models.ProtectPDFJob{
		ID:           jobID,
		UserID:       userID,
		InputFileID:  req.InputFileID,
		OutputFileID: outputFileID,
		Password:     req.Password,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.stg.Protect().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	// 4. pdfcpu orqali parol qo‘yish
	// `pdfcpu` config format: user pw + owner pw (user pw yetarli)

	args := []string{
		"encrypt",
		"-upw", req.Password,
		"-opw", req.Password,
		file.FilePath,
		outputPath,
	}

	cmd := exec.Command("pdfcpu", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.log.Error("pdf encryption failed", logger.Error(err))
		return "", err
	}

	// 5. Faylni DBga saqlash
	fi, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("output file not found", logger.Error(err))
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

	// 6. Holatni yangilash
	job.Status = "done"
	if err := s.stg.Protect().Update(ctx, job); err != nil {
		s.log.Error("failed to update job status", logger.Error(err))
		return "", err
	}

	s.log.Info("PDF protection completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *protectPDFService) GetByID(ctx context.Context, id string) (*models.ProtectPDFJob, error) {
	job, err := s.stg.Protect().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
