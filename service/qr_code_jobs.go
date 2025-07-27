package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/skip2/go-qrcode" // QR kod kutubxonasi

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type QRCodeService interface {
	Create(ctx context.Context, req models.CreateQRCodeRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.QRCodeJob, error)
}

type qrCodeService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewQRCodeService(stg storage.IStorage, log logger.ILogger) QRCodeService {
	return &qrCodeService{stg: stg, log: log}
}

func (s *qrCodeService) Create(ctx context.Context, req models.CreateQRCodeRequest, userID string) (string, error) {
	s.log.Info("QRCodeService.Create called")

	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		return "", fmt.Errorf("input file not found")
	}

	jobID := uuid.New().String()
	job := &models.QRCodeJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		QRContent:   req.QRContent,
		Position:    req.Position,
		Size:        req.Size,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.stg.QRCode().Create(ctx, job); err != nil {
		return "", err
	}

	// QR kod yaratish (PNG fayl sifatida)
	qrPngPath := filepath.Join("storage/qr_code", jobID+".png")

	err = qrcode.WriteFile(req.QRContent, qrcode.Medium, req.Size, qrPngPath)
	if err != nil {
		s.log.Error("failed to generate qr code", logger.Error(err))
		job.Status = "failed"
		_ = s.stg.QRCode().Update(ctx, job)
		return "", err
	}

	// PDF ga QR kodni joylashtirish
	outputID := uuid.New().String()
	outputPath := filepath.Join("storage/qr_code_pdf", outputID+".pdf")

	conf := model.NewDefaultConfiguration()

	// Pozitsiyani pdfcpu uchun o'zgartirish (misol uchun)
	var pos string
	switch req.Position {
	case "top-left":
		pos = "tl"
	case "top-right":
		pos = "tr"
	case "bottom-left":
		pos = "bl"
	case "bottom-right":
		pos = "br"
	case "center":
		pos = "c"
	default:
		pos = "c"
	}

	// Watermark uchun parametrlar: pozitsiya va o'lcham
	wmDetails := fmt.Sprintf("pos:%s, scale:1 rel, rot:0", pos)

	wm, err := pdfcpu.ParseImageWatermarkDetails(qrPngPath, wmDetails, true, conf.Unit)
	if err != nil {
		s.log.Error("failed to parse watermark", logger.Error(err))
		job.Status = "failed"
		_ = s.stg.QRCode().Update(ctx, job)
		return "", err
	}

	err = api.AddWatermarksFile(file.FilePath, outputPath, nil, wm, conf)
	if err != nil {
		s.log.Error("failed to add qr watermark", logger.Error(err))
		job.Status = "failed"
		_ = s.stg.QRCode().Update(ctx, job)
		return "", err
	}

	// Natija faylni DBga yozish
	fi, _ := os.Stat(outputPath)
	newFile := models.File{
		ID:         outputID,
		UserID:     &userID,
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
	err = s.stg.QRCode().Update(ctx, job)
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *qrCodeService) GetByID(ctx context.Context, id string) (*models.QRCodeJob, error) {
	return s.stg.QRCode().GetByID(ctx, id)
}
