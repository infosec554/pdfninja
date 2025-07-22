package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
)

type protectRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewProtectRepo(db *pgxpool.Pool, log logger.ILogger) *protectRepo {
	return &protectRepo{
		db:  db,
		log: log,
	}
}

func (r *protectRepo) Create(ctx context.Context, job *models.ProtectPDFJob) error {
	query := `
        INSERT INTO protect_pdf_jobs (
            id, user_id, input_file_id, output_file_id, password, status, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.Exec(ctx, query,
		job.ID,
		job.UserID,
		job.InputFileID,
		job.OutputFileID,
		job.Password,
		job.Status,
		job.CreatedAt,
	)
	return err
}

func (r *protectRepo) GetByID(ctx context.Context, id string) (*models.ProtectPDFJob, error) {
	query := `
        SELECT id, user_id, input_file_id, output_file_id, password, status, created_at
        FROM protect_pdf_jobs
        WHERE id = $1
    `
	var job models.ProtectPDFJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.OutputFileID,
		&job.Password,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *protectRepo) Update(ctx context.Context, job *models.ProtectPDFJob) error {
	query := `
        UPDATE protect_pdf_jobs
        SET status = $1
        WHERE id = $2
    `
	_, err := r.db.Exec(ctx, query, job.Status, job.ID)
	return err
}
