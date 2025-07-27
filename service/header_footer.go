package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/addheaderfooter"
	"test/pkg/logger"
	"test/storage"
)

type AddHeaderFooterService interface {
	Create(ctx context.Context, req models.CreateAddHeaderFooterRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.AddHeaderFooterJob, error)
}

type addHeaderFooterService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewAddHeaderFooterService(stg storage.IStorage, log logger.ILogger) AddHeaderFooterService {
	return &addHeaderFooterService{
		stg: stg,
		log: log,
	}
}

func (s *addHeaderFooterService) Create(ctx context.Context, req models.CreateAddHeaderFooterRequest, userID string) (string, error) {
	s.log.Info("AddHeaderFooterService.Create called")

	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		return "", fmt.Errorf("input file not found")
	}

	jobID := uuid.New().String()
	job := &models.AddHeaderFooterJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		HeaderText:  req.HeaderText,
		FooterText:  req.FooterText,
		FontSize:    req.FontSize,
		FontColor:   req.FontColor,
		Position:    req.Position,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	err = s.stg.AddHeaderFooter().Create(ctx, job)
	if err != nil {
		return "", err
	}

	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/header_footer", outputID+".pdf")

	params := addheaderfooter.AddHeaderFooterParams{
		InputPath:  file.FilePath,
		OutputPath: outputPath,
		HeaderText: req.HeaderText,
		FooterText: req.FooterText,
		PageRange:  "all", // kerak bo‘lsa o‘zgartiring
	}

	err = addheaderfooter.AddHeaderFooterPDF(params)
	if err != nil {
		s.log.Error("failed to add header/footer", logger.Error(err))
		return "", err
	}

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

	_, err = s.stg.File().Save(ctx, newFile)
	if err != nil {
		return "", err
	}

	job.OutputFileID = outputID
	job.Status = "done"
	err = s.stg.AddHeaderFooter().Update(ctx, job)
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *addHeaderFooterService) GetByID(ctx context.Context, id string) (*models.AddHeaderFooterJob, error) {
	return s.stg.AddHeaderFooter().GetByID(ctx, id)
}
