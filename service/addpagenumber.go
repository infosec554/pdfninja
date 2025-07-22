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

type AddPageNumberService interface {
	Create(ctx context.Context, req models.AddPageNumbersRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.AddPageNumberJob, error)
}

type addPageNumberService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewAddPageNumberService(stg storage.IStorage, log logger.ILogger) AddPageNumberService {
	return &addPageNumberService{
		stg: stg,
		log: log,
	}
}

func (s *addPageNumberService) Create(ctx context.Context, req models.AddPageNumbersRequest, userID string) (string, error) {
	s.log.Info("AddPageNumberService.Create called")

	// 1. Kiruvchi faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", fmt.Errorf("input file not found: %v", err)
	}

	// 2. Job yaratish
	jobID := uuid.NewString()
	outputFileID := uuid.NewString()
	outputPath := filepath.Join("storage/add_page_numbers", outputFileID+".pdf")

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		s.log.Error("failed to create output directory", logger.Error(err))
		return "", err
	}

	job := &models.AddPageNumberJob{
		ID:           jobID,
		UserID:       userID,
		InputFileID:  req.InputFileID,
		OutputFileID: outputFileID,
		Position:     req.Position,
		FirstNumber:  req.FirstNumber,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.stg.AddPageNumber().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	// 3. CLI orqali sahifalarga raqam qo‘shish (pdfcpu ishlatilmoqda)
	args := []string{
		"stamp",
		"add",
		"-mode", "text",
		"-pages", "all",
		"-pos", req.Position,
		"-text", "#p", // `#p` sahifa raqami degani
		file.FilePath,
		outputPath,
	}
	cmd := exec.Command("pdfcpu", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.log.Error("pdfcpu failed", logger.Error(err))
		return "", err
	}

	// 4. Yangi faylni DBga saqlash
	fi, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("failed to stat output file", logger.Error(err))
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

	// 5. Job holatini "done" qilish
	job.Status = "done"
	if err := s.stg.AddPageNumber().Update(ctx, job); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("Add page number job completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *addPageNumberService) GetByID(ctx context.Context, id string) (*models.AddPageNumberJob, error) {
	job, err := s.stg.AddPageNumber().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
