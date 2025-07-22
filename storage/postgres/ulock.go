package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type unlockRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUnlockRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUnlockPDFStorage {
	return &unlockRepo{db: db, log: log}
}

func (r *unlockRepo) Create(ctx context.Context, job *models.UnlockPDFJob) error {
	query := `
		INSERT INTO unlock_pdf_jobs (
			id, user_id, input_file_id, output_file_id, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID,
		job.OutputFileID, job.Status, job.CreatedAt)

	if err != nil {
		return fmt.Errorf("unlockRepo.Create: %w", err)
	}

	return nil
}

func (r *unlockRepo) GetByID(ctx context.Context, id string) (*models.UnlockPDFJob, error) {
	query := `
		SELECT id, user_id, input_file_id, output_file_id, status, created_at
		FROM unlock_pdf_jobs
		WHERE id = $1;
	`

	row := r.db.QueryRow(ctx, query, id)

	var job models.UnlockPDFJob
	err := row.Scan(
		&job.ID, &job.UserID, &job.InputFileID,
		&job.OutputFileID, &job.Status, &job.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unlockRepo.GetByID: not found")
		}
		return nil, fmt.Errorf("unlockRepo.GetByID: %w", err)
	}

	return &job, nil
}

func (r *unlockRepo) Update(ctx context.Context, job *models.UnlockPDFJob) error {
	query := `
		UPDATE unlock_pdf_jobs
		SET output_file_id = $1,
			status = $2
		WHERE id = $3;
	`

	_, err := r.db.Exec(ctx, query,
		job.OutputFileID, job.Status, job.ID)

	if err != nil {
		return fmt.Errorf("unlockRepo.Update: %w", err)
	}

	return nil
}
