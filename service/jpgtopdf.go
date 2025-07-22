package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type JpgToPdfService interface {
	Create(ctx context.Context, req models.CreateJpgToPdfRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.JpgToPdfJob, error)
}

type jpgToPdfService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewJpgToPdfService(stg storage.IStorage, log logger.ILogger) JpgToPdfService {
	return &jpgToPdfService{stg: stg, log: log}
}

func (s *jpgToPdfService) Create(ctx context.Context, req models.CreateJpgToPdfRequest, userID string) (string, error) {
	s.log.Info("JpgToPdfService.Create called")

	var inputPaths []string
	for _, imageID := range req.ImageFileIDs {
		file, err := s.stg.File().GetByID(ctx, imageID)
		if err != nil {
			s.log.Error("image file not found", logger.String("fileID", imageID), logger.Error(err))
			return "", fmt.Errorf("file not found: %s", imageID)
		}
		inputPaths = append(inputPaths, file.FilePath)
	}

	jobID := uuid.New().String()
	job := &models.JpgToPdfJob{
		ID:           jobID,
		UserID:       userID,
		ImageFileIDs: req.ImageFileIDs,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.stg.JpgToPdf().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/jpg_to_pdf", outputID+".pdf")

	// === GoPDF bilan PDF yaratish ===
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, imgPath := range inputPaths {
		pdf.AddPage()
		// A4 o‘lchamdagi rasmni joylashtirish
		pdf.ImageOptions(imgPath, 10, 10, 190, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		s.log.Error("failed to generate PDF", logger.Error(err))
		return "", err
	}

	fi, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("cannot stat output PDF", logger.Error(err))
		return "", err
	}

	newFile := models.File{
		ID:         outputID,
		UserID:     userID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   fi.Size(),
		UploadedAt: time.Now(),
	}

	if _, err := s.stg.File().Save(ctx, newFile); err != nil {
		s.log.Error("failed to save file", logger.Error(err))
		return "", err
	}

	job.OutputFileID = outputID
	job.Status = "done"
	if err := s.stg.JpgToPdf().Update(ctx, job); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("jpg to pdf job completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *jpgToPdfService) GetByID(ctx context.Context, id string) (*models.JpgToPdfJob, error) {
	job, err := s.stg.JpgToPdf().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get jpg to pdf job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
