package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/addbackground"
	"test/pkg/logger"
	"test/storage"
)

type AddBackgroundService interface {
	Create(ctx context.Context, req models.CreateAddBackgroundRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.AddBackgroundJob, error)
}

type addBackgroundService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewAddBackgroundService(stg storage.IStorage, log logger.ILogger) AddBackgroundService {
	return &addBackgroundService{
		stg: stg,
		log: log,
	}
}

func (s *addBackgroundService) Create(ctx context.Context, req models.CreateAddBackgroundRequest, userID string) (string, error) {
	s.log.Info("AddBackgroundService.Create called")

	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		return "", fmt.Errorf("input file not found")
	}

	bgFile, err := s.stg.File().GetByID(ctx, req.BackgroundImageFileID)
	if err != nil {
		return "", fmt.Errorf("background image file not found")
	}

	jobID := uuid.New().String()
	job := &models.AddBackgroundJob{
		ID:                    jobID,
		UserID:                userID,
		InputFileID:           req.InputFileID,
		BackgroundImageFileID: req.BackgroundImageFileID,
		Opacity:               req.Opacity,
		Position:              req.Position,
		Status:                "pending",
		CreatedAt:             time.Now(),
	}

	err = s.stg.AddBackground().Create(ctx, job)
	if err != nil {
		return "", err
	}

	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/add_background", outputID+".pdf")

	params := addbackground.AddBackgroundParams{
		InputPath:           file.FilePath,
		OutputPath:          outputPath,
		BackgroundImagePath: bgFile.FilePath,
		Opacity:             req.Opacity,
		Position:            req.Position,
		PageRange:           "all",
	}

	err = addbackground.AddBackgroundImage(params)
	if err != nil {
		s.log.Error("failed to add background", logger.Error(err))
		return "", err
	}

	fi, err := os.Stat(outputPath)
	if err != nil {
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

	_, err = s.stg.File().Save(ctx, newFile)
	if err != nil {
		return "", err
	}

	job.OutputFileID = outputID
	job.Status = "done"
	err = s.stg.AddBackground().Update(ctx, job)
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *addBackgroundService) GetByID(ctx context.Context, id string) (*models.AddBackgroundJob, error) {
	return s.stg.AddBackground().GetByID(ctx, id)
}
