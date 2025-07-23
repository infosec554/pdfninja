package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type ExtractPageService interface {
	Create(ctx context.Context, req models.ExtractPagesRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.ExtractJob, error)
}

type extractPageService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewExtractService(stg storage.IStorage, log logger.ILogger) ExtractPageService {
	return &extractPageService{
		stg: stg,
		log: log,
	}
}

func (s *extractPageService) Create(ctx context.Context, req models.ExtractPagesRequest, userID string) (string, error) {
	s.log.Info("ExtractService.Create called")

	// 1. Kirish faylini olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	jobID := uuid.New().String()
	job := &models.ExtractJob{
		ID:             jobID,
		UserID:         userID,
		InputFileID:    req.InputFileID,
		PagesToExtract: req.PageRanges,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	// 3. DBga yozish
	if err := s.stg.ExtractPage().Create(ctx, job); err != nil {
		s.log.Error("failed to create extract job", logger.Error(err))
		return "", err
	}

	// 4. Fayllarni ajratish
	inputPath := file.FilePath
	outputDir := filepath.Join("storage/extract", job.ID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output dir", logger.Error(err))
		return "", err
	}

	ranges := strings.Split(req.PageRanges, ",") // masalan: "1-2,4,6"
	config := model.NewDefaultConfiguration()

	var firstOutputFileID *string

	for i, r := range ranges {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("range_%d.pdf", i+1))

		// Sahifalarni ajratish
		err := api.ExtractPagesFile(inputPath, outputPath, []string{r}, config)
		if err != nil {
			s.log.Error("failed to extract range", logger.Error(err))
			continue
		}

		// Faylni tekshirish
		fi, err := os.Stat(outputPath)
		if err != nil {
			s.log.Error("failed to stat extracted file", logger.Error(err))
			continue
		}

		// Faylni saqlash
		newFile := models.File{
			ID:         uuid.New().String(),
			UserID:     userID,
			FileName:   filepath.Base(outputPath),
			FilePath:   outputPath,
			FileType:   "application/pdf",
			FileSize:   fi.Size(),
			UploadedAt: time.Now(),
		}

		fileID, err := s.stg.File().Save(ctx, newFile)
		if err != nil {
			s.log.Error("failed to save extracted file", logger.Error(err))
			continue
		}

		// Faqat 1-chi faylni output_file_id sifatida saqlaymiz
		if firstOutputFileID == nil {
			firstOutputFileID = &fileID
		}
	}

	// 5. Job holatini yangilash
	job.Status = "done"
	job.OutputFileID = firstOutputFileID

	if err := s.stg.ExtractPage().Update(ctx, job); err != nil {
		s.log.Error("failed to update extract job", logger.Error(err))
		return "", err
	}

	s.log.Info("extract job completed", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *extractPageService) GetByID(ctx context.Context, id string) (*models.ExtractJob, error) {
	job, err := s.stg.ExtractPage().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get extract job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
