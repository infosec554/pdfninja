package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/createpdffortransalatepdf" // Fayldan matn chiqarish uchun util
	"test/pkg/logger"
	"test/storage"
)

type PDFTextSearchService interface {
	Create(ctx context.Context, inputFileID  string, userID *string) (string, error)
	GetByID(ctx context.Context, id string) (*models.PDFTextSearchJob, error)
}

type pdfTextSearchService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewPDFTextSearchService(stg storage.IStorage, log logger.ILogger) PDFTextSearchService {
	return &pdfTextSearchService{stg: stg, log: log}
}

func (s *pdfTextSearchService) Create(ctx context.Context, inputFileID string, userID *string) (string, error) {
	s.log.Info("PDFTextSearchService.Create called")

	file, err := s.stg.File().GetByID(ctx, inputFileID)
	if err != nil {
		return "", fmt.Errorf("file not found")
	}

	jobID := uuid.New().String()
	job := &models.PDFTextSearchJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: inputFileID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.PDFTextSearch().Create(ctx, job); err != nil {
		return "", err
	}

	// Fayldan matnni olish (extract)
	text, err := createpdffortransalatepdf.ExtractTextFromPDF(file.FilePath)
	if err != nil {
		s.log.Error("failed to extract text", logger.Error(err))
		job.Status = "failed"
		_ = s.stg.PDFTextSearch().Update(ctx, job)
		return "", err
	}

	job.ExtractedText = text
	job.Status = "done"

	if err := s.stg.PDFTextSearch().Update(ctx, job); err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *pdfTextSearchService) GetByID(ctx context.Context, id string) (*models.PDFTextSearchJob, error) {
	return s.stg.PDFTextSearch().GetByID(ctx, id)
}
