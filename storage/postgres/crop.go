package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

// cropRepo implements ICropPDFStorage
type cropRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

// NewCropRepo konstruktori
func NewCropRepo(db *pgxpool.Pool, log logger.ILogger) storage.ICropPDFStorage {
	return &cropRepo{db: db, log: log}
}

// Create – yangi crop job qo‘shadi
func (r *cropRepo) Create(ctx context.Context, job *models.CropPDFJob) error {
	query := `
        INSERT INTO crop_pdf_jobs
        (id, user_id, input_file_id, top, bottom, left, right, status, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    `
	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID,
		job.Top, job.Bottom, job.Left, job.Right,
		job.Status, job.CreatedAt,
	)
	if err != nil {
		r.log.Error("cropRepo.Create failed", logger.Error(err))
	}
	return err
}

// GetByID – ID bo‘yicha crop job ma’lumotini qaytaradi
func (r *cropRepo) GetByID(ctx context.Context, id string) (*models.CropPDFJob, error) {
	query := `
        SELECT id, user_id, input_file_id, top, bottom, left, right,
               output_file_id, status, created_at
        FROM crop_pdf_jobs
        WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)

	var job models.CropPDFJob
	err := row.Scan(
		&job.ID, &job.UserID, &job.InputFileID,
		&job.Top, &job.Bottom, &job.Left, &job.Right,
		&job.OutputFileID, &job.Status, &job.CreatedAt,
	)
	if err != nil {
		r.log.Error("cropRepo.GetByID failed", logger.Error(err))
		return nil, err
	}
	return &job, nil
}

// Update – crop job natijasi va statusini yangilaydi
func (r *cropRepo) Update(ctx context.Context, job *models.CropPDFJob) error {
	query := `
        UPDATE crop_pdf_jobs
        SET output_file_id = $1, status = $2
        WHERE id = $3
    `
	_, err := r.db.Exec(ctx, query,
		job.OutputFileID, job.Status, job.ID,
	)
	if err != nil {
		r.log.Error("cropRepo.Update failed", logger.Error(err))
	}
	return err
}
