package service

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/gotenberg"
	"test/pkg/logger"
	"test/storage"
)

type HTMLToPDFService interface {
	Create(ctx context.Context, req models.CreateHTMLToPDFRequest, userID *string) (string, error)
	GetByID(ctx context.Context, id string) (*models.HTMLToPDFJob, error)
}

type htmlToPDFService struct {
	stg       storage.IStorage
	log       logger.ILogger
	gotClient gotenberg.Client
}

func NewHTMLToPDFService(stg storage.IStorage, log logger.ILogger, gotClient gotenberg.Client) HTMLToPDFService {
	return &htmlToPDFService{
		stg:       stg,
		log:       log,
		gotClient: gotClient,
	}
}

func (s *htmlToPDFService) Create(ctx context.Context, req models.CreateHTMLToPDFRequest, userID *string) (string, error) {
	s.log.Info("HTMLToPDFService.Create called")

	// 1. HTML faylni vaqtinchalik yaratish
	tempID := uuid.NewString()
	htmlFilePath := filepath.Join("tmp", tempID+".html")

	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		s.log.Error("failed to create tmp dir", logger.Error(err))
		return "", err
	}

	if err := os.WriteFile(htmlFilePath, []byte(req.HTMLContent), 0644); err != nil {
		s.log.Error("failed to write html file", logger.Error(err))
		return "", err
	}

	// 2. Gotenberg orqali PDFga aylantirish
	resultBytes, err := s.gotClient.HTMLToPDF(ctx, htmlFilePath)
	if err != nil {
		s.log.Error("Gotenberg HTML conversion failed", logger.Error(err))
		return "", err
	}

	// 3. Yangi ID va chiqish faylini saqlash
	jobID := uuid.NewString()
	outputFileID := uuid.NewString()
	outputPath := filepath.Join("storage/html_to_pdf", outputFileID+".pdf")

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return "", err
	}
	if err := os.WriteFile(outputPath, resultBytes, 0644); err != nil {
		return "", err
	}

	fi, _ := os.Stat(outputPath)
	newFile := models.File{
		ID:         outputFileID,
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

	// 4. Jobni saqlash
	job := &models.HTMLToPDFJob{
		ID:           jobID,
		UserID:       userID,
		OutputFileID: &outputFileID,
		Status:       "done",
		CreatedAt:    time.Now(),
	}
	if err := s.stg.HTMLToPDF().Create(ctx, job); err != nil {
		s.log.Error("failed to save html-to-pdf job", logger.Error(err))
		return "", err
	}

	// 5. Vaqtinchalik HTML faylni o‘chirish
	_ = os.Remove(htmlFilePath)

	return jobID, nil
}

func (s *htmlToPDFService) GetByID(ctx context.Context, id string) (*models.HTMLToPDFJob, error) {
	job, err := s.stg.HTMLToPDF().GetByID(ctx, id)
	if err != nil {
		s.log.Error("HTMLToPDF job not found", logger.Error(err))
		return nil, err
	}
	return job, nil
}
