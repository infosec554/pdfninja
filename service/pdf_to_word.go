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

type PdfToWordService interface {
	Create(ctx context.Context, req models.PDFToWordRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.PDFToWordJob, error)
}

type pdfToWordService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewPdfToWordService(stg storage.IStorage, log logger.ILogger) PdfToWordService {
	return &pdfToWordService{
		stg: stg,
		log: log,
	}
}

func (s *pdfToWordService) Create(ctx context.Context, req models.PDFToWordRequest, userID string) (string, error) {
	s.log.Info("PdfToWordService.Create called")

	// 1. Input faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("failed to get input PDF file", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	jobID := uuid.NewString()
	outputDir := filepath.Join("storage/pdf_to_word")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output directory", logger.Error(err))
		return "", err
	}

	outputPath := filepath.Join(outputDir, jobID+".docx")

	job := &models.PDFToWordJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		OutputPath:  outputPath,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.PDFToWord().Create(ctx, job); err != nil {
		s.log.Error("failed to store pdf_to_word job", logger.Error(err))
		return "", err
	}

	// 3. LibreOffice orqali convert qilish
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "docx", "--outdir", outputDir, file.FilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.log.Error("failed to convert PDF to Word", logger.Error(err))
		return "", err
	}

	// 4. Holatni yangilash
	job.Status = "done"
	if err := s.stg.PDFToWord().Update(ctx, job); err != nil {
		s.log.Error("failed to update job status", logger.Error(err))
		return "", err
	}

	s.log.Info("PDF to Word conversion completed", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *pdfToWordService) GetByID(ctx context.Context, id string) (*models.PDFToWordJob, error) {
	job, err := s.stg.PDFToWord().GetByID(ctx, id) // ✅ to‘g‘ri storage chaqirildi
	if err != nil {
		s.log.Error("failed to get PDFToWordJob", logger.Error(err))
		return nil, err
	}
	return job, nil
}
