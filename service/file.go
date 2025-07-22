package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type FileService interface {
	Upload(ctx context.Context, req models.File) (string, error)
	Get(ctx context.Context, id string) (models.File, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]models.File, error)
}

type fileService struct {
	stg storage.IFileStorage
	log logger.ILogger
}

func NewFileService(stg storage.IStorage, log logger.ILogger) FileService {
	return &fileService{
		stg: stg.File(), // storageManager.File()
		log: log,
	}
}

// Upload - faylni DBga yozish (upload jarayoni)
func (s *fileService) Upload(ctx context.Context, req models.File) (string, error) {
	s.log.Info("FileService.Upload called", logger.String("file_name", req.FileName))

	req.ID = uuid.NewString()
	req.UploadedAt = time.Now()

	id, err := s.stg.Save(ctx, req)
	if err != nil {
		s.log.Error("failed to save file", logger.Error(err))
		return "", err
	}
	return id, nil
}

// Get - faylni ID orqali olish
func (s *fileService) Get(ctx context.Context, id string) (models.File, error) {
	return s.stg.GetByID(ctx, id)
}

// Delete - faylni o‘chirish
func (s *fileService) Delete(ctx context.Context, id string) error {
	return s.stg.Delete(ctx, id)
}

// List - user fayllari ro‘yxati
func (s *fileService) List(ctx context.Context, userID string) ([]models.File, error) {
	return s.stg.ListByUser(ctx, userID)
}
