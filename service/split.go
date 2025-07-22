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

type SplitService interface {
	Create(ctx context.Context, req models.CreateSplitJobRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.SplitJob, error)
}

type splitService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewSplitService(stg storage.IStorage, log logger.ILogger) SplitService {
	return &splitService{
		stg: stg,
		log: log,
	}
}

func (s *splitService) Create(ctx context.Context, req models.CreateSplitJobRequest, userID string) (string, error) {
	s.log.Info("SplitService.Create called")

	// 1. Faylni bazadan olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("input file not found", logger.Error(err))
		return "", err
	}

	// 2. SplitJob struct
	job := &models.SplitJob{
		ID:            uuid.New().String(),
		UserID:        userID,
		InputFileID:   req.InputFileID,
		SplitRanges:   req.SplitRanges,
		OutputFileIDs: []string{},
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	// 3. Jobni bazaga yozish
	if err := s.stg.Split().Create(ctx, job); err != nil {
		s.log.Error("failed to create split job", logger.Error(err))
		return "", err
	}

	// 4. Split qilish
	inputPath := file.FilePath
	outputDir := fmt.Sprintf("storage/split/%s", job.ID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output dir", logger.Error(err))
		return "", err
	}

	config := model.NewDefaultConfiguration()

	ranges := strings.Split(req.SplitRanges, ",") // masalan: "1-3,4-6"
	for i, r := range ranges {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("part_%d.pdf", i+1))
		err := api.ExtractPagesFile(inputPath, outputPath, []string{r}, config)
		if err != nil {
			s.log.Error("failed to extract range", logger.Error(err))
			continue
		}

		// Fayl haqida ma’lumot
		fi, err := os.Stat(outputPath)
		if err != nil {
			s.log.Error("failed to stat output file", logger.Error(err))
			continue
		}

		outputFile := models.File{
			ID:         uuid.New().String(),
			UserID:     userID,
			FileName:   filepath.Base(outputPath),
			FilePath:   outputPath,
			FileType:   "application/pdf",
			FileSize:   fi.Size(),
			UploadedAt: time.Now(),
		}

		fileID, err := s.stg.File().Save(ctx, outputFile)
		if err != nil {
			s.log.Error("failed to save output file", logger.Error(err))
			continue
		}

		job.OutputFileIDs = append(job.OutputFileIDs, fileID)
	}

	// 5. Jobni yangilash
	if err := s.stg.Split().UpdateOutputFiles(ctx, job.ID, job.OutputFileIDs); err != nil {
		s.log.Error("failed to update split job", logger.Error(err))
		return "", err
	}

	s.log.Info("split job completed", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *splitService) GetByID(ctx context.Context, id string) (*models.SplitJob, error) {
	job, err := s.stg.Split().GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get split job", logger.Error(err))
		return nil, err
	}
	return job, nil
}
