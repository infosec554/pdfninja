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

type PDFToJPGService interface {
	Create(ctx context.Context, req models.PDFToJPGRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error)
}

type pdfToJPGService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewPDFToJPGService(stg storage.IStorage, log logger.ILogger) PDFToJPGService {
	return &pdfToJPGService{stg: stg, log: log}
}

func (s *pdfToJPGService) Create(ctx context.Context, req models.PDFToJPGRequest, userID string) (string, error) {
	s.log.Info("PDFToJPGService.Create called")

	// 1. PDF faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("PDF file not found", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	jobID := uuid.NewString()
	outputDir := filepath.Join("storage/pdf_to_jpg", jobID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output dir", logger.Error(err))
		return "", err
	}

	job := &models.PDFToJPGJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.PDFToJPG().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	// 3. CLI orqali sahifalarni JPG formatda chiqarish (pdftoppm kerak)
	// 📌 install: sudo apt install poppler-utils
	outputPath := filepath.Join(outputDir, "page")
	cmd := exec.Command("pdftoppm", "-jpeg", file.FilePath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.log.Error("failed to convert PDF to JPG", logger.Error(err))
		return "", err
	}

	// 4. Yaratilgan JPG fayllarni to‘plash
	var outputFiles []string
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			outputFiles = append(outputFiles, path)
		}
		return nil
	})
	if err != nil {
		s.log.Error("error collecting jpg files", logger.Error(err))
		return "", err
	}

	job.OutputPaths = outputFiles
	job.Status = "done"
	if err := s.stg.PDFToJPG().Update(ctx, job); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("PDF to JPG completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *pdfToJPGService) GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error) {
	job, err := s.stg.PDFToJPG().GetByID(ctx, id)
	if err != nil {
		s.log.Error("job not found", logger.Error(err))
		return nil, err
	}
	return job, nil
}
