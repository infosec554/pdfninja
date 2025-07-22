package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type MergeService interface {
	Create(ctx context.Context, userID string, inputFileIDs []string) (string, error)
	GetByID(ctx context.Context, id string) (*models.MergeJob, error)
	ProcessJob(ctx context.Context, jobID string) (string, error)
}

type mergeService struct {
	stg storage.IMergeStorage
	fs  storage.IFileStorage
	log logger.ILogger
}

func NewMergeService(stg storage.IStorage, log logger.ILogger) MergeService {
	return &mergeService{
		stg: stg.Merge(),
		fs:  stg.File(),
		log: log,
	}
}

func (s *mergeService) Create(ctx context.Context, userID string, inputFileIDs []string) (string, error) {
	s.log.Info("MergeService.Create called", logger.String("userID", userID))

	job := &models.MergeJob{
		ID:        uuid.New().String(),
		UserID:    userID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := s.stg.Create(ctx, job); err != nil {
		s.log.Error("failed to create merge job", logger.Error(err))
		return "", err
	}

	if err := s.stg.AddInputFiles(ctx, job.ID, inputFileIDs); err != nil {
		s.log.Error("failed to add input files", logger.Error(err))
		return "", err
	}

	s.log.Info("merge job created", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *mergeService) GetByID(ctx context.Context, id string) (*models.MergeJob, error) {
	s.log.Info("MergeService.GetByID called", logger.String("jobID", id))

	job, err := s.stg.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get merge job", logger.Error(err))
		return nil, err
	}

	return job, nil
}

func (s *mergeService) ProcessJob(ctx context.Context, jobID string) (string, error) {
	s.log.Info("📥 MergeService.ProcessJob called", logger.String("jobID", jobID))

	// Merge katalogini tayyorlash
	if err := os.MkdirAll("storage/merge", os.ModePerm); err != nil {
		s.log.Error("❌ failed to create merge dir", logger.Error(err))
		return "", err
	}

	// Job ni olish
	job, err := s.stg.GetByID(ctx, jobID)
	if err != nil {
		s.log.Error("❌ job not found", logger.Error(err))
		return "", fmt.Errorf("job not found: %w", err)
	}

	if len(job.InputFileIDs) < 2 {
		s.log.Error("❌ not enough input files", logger.Any("input_ids", job.InputFileIDs))
		return "", fmt.Errorf("need at least two files to merge")
	}

	// Fayl yo'llarini tayyorlash
	var inputPaths []string
	for _, fileID := range job.InputFileIDs {
		file, err := s.fs.GetByID(ctx, fileID)
		if err != nil {
			s.log.Error("❌ failed to get input file", logger.String("fileID", fileID), logger.Error(err))
			return "", fmt.Errorf("failed to get file: %w", err)
		}
		s.log.Info("✅ got input file", logger.String("fileID", fileID), logger.String("path", file.FilePath))

		// Fayl mavjudligini tekshiramiz
		if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
			s.log.Error("❌ input file does not exist", logger.String("filePath", file.FilePath))
			return "", fmt.Errorf("input file does not exist: %s", file.FilePath)
		}

		inputPaths = append(inputPaths, file.FilePath)
	}

	s.log.Info("📄 All input files ready", logger.Any("paths", inputPaths))

	// Output fayl tayyorlash
	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/merge", outputID+".pdf")

	config := model.NewDefaultConfiguration()
	if err := api.MergeCreateFile(inputPaths, outputPath, false, config); err != nil {
		s.log.Error("❌ pdf merge failed", logger.Error(err))
		return "", fmt.Errorf("merge failed: %w", err)
	}
	s.log.Info("✅ PDF successfully merged", logger.String("outputPath", outputPath))

	// Output faylni tekshirish
	info, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("❌ failed to stat merged file", logger.Error(err))
		return "", fmt.Errorf("output file stat failed: %w", err)
	}

	// Yangi faylni saqlash
	outputFile := models.File{
		ID:         outputID,
		UserID:     job.UserID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   info.Size(),
		UploadedAt: time.Now(),
	}
	if _, err = s.fs.Save(ctx, outputFile); err != nil {
		s.log.Error("❌ failed to save merged file", logger.Error(err))
		return "", fmt.Errorf("failed to save output file: %w", err)
	}

	// Jobni yangilash
	job.OutputFileID = &outputID
	job.Status = "done"
	if err := s.stg.Update(ctx, job); err != nil {
		s.log.Error("❌ failed to update merge job", logger.Error(err))
		return "", fmt.Errorf("failed to update merge job: %w", err)
	}

	s.log.Info("🎉 Merge job done", logger.String("outputID", outputID))
	return outputID, nil
}
