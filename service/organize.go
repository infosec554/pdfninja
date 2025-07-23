package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type OrganizeService interface {
	Create(ctx context.Context, req models.CreateOrganizeJobRequest, userID string) (string, error)
	GetByID(ctx context.Context, id string) (*models.OrganizeJob, error)
}

type organizeService struct {
	stg storage.IStorage
	log logger.ILogger
}

func NewOrganizeService(stg storage.IStorage, log logger.ILogger) OrganizeService {
	return &organizeService{
		stg: stg,
		log: log,
	}
}

func (s *organizeService) Create(ctx context.Context, req models.CreateOrganizeJobRequest, userID string) (string, error) {
	s.log.Info("OrganizeService.Create called")

	// 1. Faylni olish
	file, err := s.stg.File().GetByID(ctx, req.InputFileID)
	if err != nil {
		s.log.Error("Input file not found", logger.Error(err))
		return "", err
	}

	// 2. Organize job obyektini yaratish
	jobID := uuid.New().String()
	job := &models.OrganizeJob{
		ID:          jobID,
		UserID:      userID,
		InputFileID: req.InputFileID,
		NewOrder:    req.NewOrder,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	s.log.Info("Creating new organize job", logger.Any("job", job))

	// 3. Organize job-ni DB ga saqlash
	if err := s.stg.Organize().Create(ctx, job); err != nil {
		s.log.Error("Failed to create organize job in DB", logger.Error(err))
		return "", err
	}

	// 4. Sahifa tartibini []int → "3,1,2" ga aylantirish
	orderList, err := parsePageOrder(req.NewOrder)
	if err != nil {
		s.log.Error("Invalid page order format", logger.Error(err))
		return "", err
	}
	var orderStr []string
	for _, page := range orderList {
		orderStr = append(orderStr, strconv.Itoa(page))
	}
	finalOrder := strings.Join(orderStr, ",")

	// 5. Output fayl uchun path va katalog tayyorlash
	outputID := uuid.New().String()
	outputDir := "storage/organize"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		s.log.Error("Failed to create output directory", logger.Error(err))
		return "", err
	}
	outputPath := filepath.Join(outputDir, outputID+".pdf")

	// 6. pdfcpu reorder buyrug'ini bajarish
cmd := exec.Command("pdfcpu", "reorder", "-pages", finalOrder, file.FilePath, outputPath)
output, err := cmd.CombinedOutput()
if err != nil {
	s.log.Error("pdfcpu execution failed", logger.String("stderr", string(output)), logger.Error(err))
	return "", fmt.Errorf("pdfcpu error: %s", string(output))
}

	// 7. Output faylni diskdan olish
	fi, err := os.Stat(outputPath)
	if err != nil {
		s.log.Error("Output file stat failed", logger.Error(err))
		return "", err
	}

	// 8. Output faylni DB ga saqlash
	newFile := models.File{
		ID:         outputID,
		UserID:     userID,
		FileName:   filepath.Base(outputPath),
		FilePath:   outputPath,
		FileType:   "application/pdf",
		FileSize:   fi.Size(),
		UploadedAt: time.Now(),
	}
	if _, err := s.stg.File().Save(ctx, newFile); err != nil {
		s.log.Error("Saving output file failed", logger.Error(err))
		return "", err
	}

	// 9. Job statusni yangilash
	job.OutputFileID = outputID
	job.Status = "done"
	if err := s.stg.Organize().Update(ctx, job); err != nil {
		s.log.Error("Job update failed", logger.Error(err))
		return "", err
	}

	s.log.Info("Organize job completed successfully", logger.String("jobID", job.ID))
	return job.ID, nil
}

func (s *organizeService) GetByID(ctx context.Context, id string) (*models.OrganizeJob, error) {
	job, err := s.stg.Organize().GetByID(ctx, id)
	if err != nil {
		s.log.Error("GetByID failed", logger.Error(err))
		return nil, err
	}
	return job, nil
}

// parsePageOrder - string tartibni int listga aylantiradi: "3,1,2" => []int{3,1,2}
func parsePageOrder(order []int) ([]int, error) {
	for _, page := range order {
		if page <= 0 {
			return nil, fmt.Errorf("invalid page number: %d", page)
		}
	}
	return order, nil
}
