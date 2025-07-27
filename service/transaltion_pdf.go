package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/createpdffortransalatepdf"
	"test/pkg/logger"
	"test/storage"
)

type TranslatePDFService interface {
	Create(ctx context.Context, req models.TranslatePDFRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.TranslatePDFJob, error)
}

type translatePDFService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewTranslatePDFService(stg storage.IStorage, log logger.ILogger) TranslatePDFService {
	return &translatePDFService{stg: stg, log: log}
}

func (s *translatePDFService) Create(ctx context.Context, req models.TranslatePDFRequest, userID string) (string, error) {
	s.log.Info("TranslatePDFService.Create called")

	// 1. Input file
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		return "", fmt.Errorf("file not found")
	}

	jobID := uuid.New().String()
	job := &models.TranslatePDFJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		TargetLang:  req.TargetLang,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.TranslatePDF().Create(ctx, job); err != nil {
		return "", err
	}

	text, err := createpdffortransalatepdf.ExtractTextFromPDF(file.FilePath)
	if err != nil {
		s.log.Error("failed to extract text", logger.Error(err))
		return "", err
	}

	// 3. Tarjima qilish (Google Translate API)
	translated, err := createpdffortransalatepdf.TranslateGoogleAPI(text, req.TargetLang)
	if err != nil {
		s.log.Error("translation failed", logger.Error(err))
		return "", err
	}

	// 4. PDF yaratish
	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/translated_pdf", outputID+".pdf")
	if err := createpdffortransalatepdf.CreatePDF(translated, outputPath); err != nil {
		return "", err
	}

	// 5. Faylni DBga yozish
	fi, _ := os.Stat(outputPath)
	newFile := models.File{
		ID:         outputID,
		UserID:     &userID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   fi.Size(),
		UploadedAt: time.Now(),
	}
	if _, err := s.stg.File().Save(ctx, newFile); err != nil {
		return "", err
	}

	job.OutputFileID = outputID
	job.Status = "done"
	if err := s.stg.TranslatePDF().Update(ctx, job); err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *translatePDFService) GetByID(ctx context.Context, id string) (*models.TranslatePDFJob, error) {
	return s.stg.TranslatePDF().GetByID(ctx, id)
}

//***************************************
