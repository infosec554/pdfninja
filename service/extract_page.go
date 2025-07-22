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

	// 1. Kiruvchi faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", err
	}

	// 2. Job yaratish
	job := &models.ExtractJob{
		ID:            uuid.New().String(),
		UserID:        userID,
		InputFileID:   req.InputFileID,
		PageRanges:    req.PageRanges,
		OutputFileIDs: []string{},
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	// 3. Job'ni DBga saqlash
	if err := s.stg.ExtractPage().Create(ctx, job); err != nil {
		s.log.Error("failed to create extract job", logger.Error(err))
		return "", err
	}

	// 4. PDF sahifalarni ajratish
	inputPath := file.FilePath
	outputDir := filepath.Join("storage/extract", job.ID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output dir", logger.Error(err))
		return "", err
	}

	ranges := strings.Split(req.PageRanges, ",") // masalan: "1-2,4"
	config := model.NewDefaultConfiguration()

	for i, r := range ranges {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("extracted_part_%d.pdf", i+1))
		err := api.ExtractPagesFile(inputPath, outputPath, []string{r}, config)
		if err != nil {
			s.log.Error("failed to extract range", logger.Error(err))
			continue
		}

		fi, err := os.Stat(outputPath)
		if err != nil {
			s.log.Error("failed to stat extracted file", logger.Error(err))
			continue
		}

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

		job.OutputFileIDs = append(job.OutputFileIDs, fileID)
	}

	// 5. Jobni yangilash
	job.Status = "done"
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
