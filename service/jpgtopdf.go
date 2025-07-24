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

type JPGToPDFService interface {
	CreateJob(ctx context.Context, userID string, inputFileIDs []string) (string, error)
	GetJobByID(ctx context.Context, id string) (*models.JPGToPDFJob, error)
}

type jpgToPDFService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewJPGToPDFService(stg storage.IStorage, log logger.ILogger) JPGToPDFService {
	return &jpgToPDFService{stg: stg, log: log}
}
func (s *jpgToPDFService) CreateJob(ctx context.Context, userID string, inputFileIDs []string) (string, error) {
	s.log.Info("JPGToPDFService.CreateJob called")

	if len(inputFileIDs) == 0 {
		return "", fmt.Errorf("no input files provided")
	}

	var inputPaths []string
	for _, fileID := range inputFileIDs {
		file, err := s.stg.File().GetByID(ctx, fileID)
		if err != nil {
			s.log.Error("file not found", logger.String("fileID", fileID), logger.Error(err))
			return "", fmt.Errorf("file not found: %s", fileID)
		}
		inputPaths = append(inputPaths, file.FilePath)
	}

	jobID := uuid.New().String()
	job := &models.JPGToPDFJob{
		ID:           jobID,
		UserID:       userID,
		InputFileIDs: inputFileIDs,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.stg.JPGToPDF().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	// 🔧 Ensure output directory exists
	outputID := uuid.New().String()
	outputDir := filepath.Join("storage", "jpg_to_pdf")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output directory", logger.Error(err))
		return "", err
	}

	outputPath := filepath.Join(outputDir, outputID+".pdf")

	// 🖼️ Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, imgPath := range inputPaths {
		pdf.AddPage()
		pdf.ImageOptions(imgPath, 10, 10, 190, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		s.log.Error("failed to generate PDF", logger.Error(err))
		return "", err
	}

	// 📁 Save output file metadata
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
		s.log.Error("failed to save output file", logger.Error(err))
		return "", err
	}

	// 🔁 Update job status and output file
	job.OutputFileID = outputID
	job.Status = "done"

	if err := s.stg.JPGToPDF().UpdateStatusAndOutput(ctx, job.ID, job.Status, job.OutputFileID); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("JPGToPDF job completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *jpgToPDFService) GetJobByID(ctx context.Context, id string) (*models.JPGToPDFJob, error) {
	job, err := s.stg.JPGToPDF().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get JPGToPDFJob", logger.Error(err))
		return nil, err
	}
	return job, nil
}
