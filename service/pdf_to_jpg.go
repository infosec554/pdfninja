package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	createzipfromfiles "test/pkg/createZipFromFiles"
	"test/pkg/logger"
	"test/storage"
)

type PDFToJPGService interface {
	Create(ctx context.Context, req models.PDFToJPGRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error)
}

type pdfToJPGService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewPDFToJPGService(stg storage.IStorage, log logger.ILogger) PDFToJPGService {
	return &pdfToJPGService{stg: stg, log: log}
}

func (s *pdfToJPGService) Create(ctx context.Context, req models.PDFToJPGRequest, userID string) (string, error) {
	s.log.Info("PDFToJPGService.Create called")

	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("PDF file not found", logger.Error(err))
		return "", err
	}

	jobID := uuid.NewString()
	outputDir := filepath.Join("storage/pdf_to_jpg", jobID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("failed to create output dir", logger.Error(err))
		return "", err
	}

	job := &models.PDFToJPGJob{
		ID:            jobID,
		UserID:        userID,
		InputFileID:   req.InputFileID,
		OutputFileIDs: []string{},
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	if err := s.stg.PDFToJPG().Create(ctx, job); err != nil {
		s.log.Error("failed to create job", logger.Error(err))
		return "", err
	}

	// 1. PDF -> JPG
	outputPath := filepath.Join(outputDir, "page")
	cmd := exec.Command("pdftoppm", "-jpeg", file.FilePath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.log.Error("failed to convert PDF to JPG", logger.Error(err))
		return "", err
	}

	// 2. JPG fayllarni yig‘ish
	var jpgPaths []string
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			jpgPaths = append(jpgPaths, path)

			fileID := uuid.NewString()
			newFile := models.File{
				ID:         fileID,
				UserID:     userID,
				FileName:   filepath.Base(path),
				FilePath:   path,
				FileType:   "image/jpeg",
				FileSize:   info.Size(),
				UploadedAt: time.Now(),
			}
			if _, err := s.stg.File().Save(ctx, newFile); err != nil {
				s.log.Error("failed to save jpg file", logger.Error(err))
				return err
			}
			job.OutputFileIDs = append(job.OutputFileIDs, fileID)
		}
		return nil
	})
	if err != nil {
		s.log.Error("error collecting jpg files", logger.Error(err))
		return "", err
	}

	// 3. ZIP yaratish
	zipPath := filepath.Join("storage/pdf_to_jpg", jobID+".zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		s.log.Error("failed to create zip file", logger.Error(err))
		return "", err
	}
	defer zipFile.Close()

	if err := createzipfromfiles.CreateZipFromFiles(zipFile, jpgPaths); err != nil {
		s.log.Error("failed to write zip", logger.Error(err))
		return "", err
	}

	// 4. ZIP faylni bazaga saqlash
	info, err := os.Stat(zipPath)
	if err != nil {
		s.log.Error("stat zip failed", logger.Error(err))
		return "", err
	}

	zipID := uuid.NewString()
	zipModel := models.File{
		ID:         zipID,
		UserID:     userID,
		FileName:   filepath.Base(zipPath),
		FilePath:   zipPath,
		FileType:   "application/zip",
		FileSize:   info.Size(),
		UploadedAt: time.Now(),
	}
	if _, err := s.stg.File().Save(ctx, zipModel); err != nil {
		s.log.Error("failed to save zip file", logger.Error(err))
		return "", err
	}

	// 5. Job holatini yangilash
	job.ZipFileID = &zipID
	job.Status = "done"
	if err := s.stg.PDFToJPG().Update(ctx, job); err != nil {
		s.log.Error("failed to update job", logger.Error(err))
		return "", err
	}

	s.log.Info("PDF to JPG ZIP completed", logger.String("jobID", jobID))
	return jobID, nil
}

func (s *pdfToJPGService) GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error) {
	job, err := s.stg.PDFToJPG().GetByID(ctx, id)
	if err != nil {
		s.log.Error("job not found", logger.Error(err))
		return nil, err
	}
	return job, nil
}
