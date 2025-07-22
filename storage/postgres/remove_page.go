package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type removeRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewRemovePageRepo(db *pgxpool.Pool, log logger.ILogger) storage.IRemovePageStorage {
	return &removeRepo{
		db:  db,
		log: log,
	}
}

func (r *removeRepo) Create(ctx context.Context, job *models.RemoveJob) error {
	query := `
		INSERT INTO remove_jobs (id, user_id, input_file_id, pages_to_remove, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID, job.PagesToRemove, job.Status, job.CreatedAt)
	return err
}

func (r *removeRepo) Update(ctx context.Context, job *models.RemoveJob) error {
	query := `
		UPDATE remove_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	return err
}

func (r *removeRepo) GetByID(ctx context.Context, id string) (*models.RemoveJob, error) {
	query := `
		SELECT id, user_id, input_file_id, pages_to_remove, output_file_id, status, created_at
		FROM remove_jobs
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var job models.RemoveJob
	err := row.Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.PagesToRemove,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}
