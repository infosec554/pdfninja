package service

import (
	"context"
	"os"
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
	CleanupOldFiles(ctx context.Context, olderThanDays int) (int, error)
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
func (s *fileService) CleanupOldFiles(ctx context.Context, olderThanDays int) (int, error) {
	s.log.Info("Cleaning up old files...")

	oldFiles, err := s.stg.GetOldFiles(ctx, olderThanDays)
	if err != nil {
		s.log.Error("failed to get old files", logger.Error(err))
		return 0, err
	}

	count := 0
	for _, file := range oldFiles {
		if err := os.Remove(file.FilePath); err != nil {
			s.log.Error("failed to delete file from disk", logger.Error(err), logger.String("path", file.FilePath))
			continue
		}
		if err := s.stg.DeleteByID(ctx, file.ID); err != nil {
			s.log.Error("failed to delete file from db", logger.Error(err), logger.String("id", file.ID))
			continue
		}
		count++
	}

	return count, nil
}
