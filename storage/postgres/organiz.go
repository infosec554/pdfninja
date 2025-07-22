package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type organizeRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewOrganizeRepo(db *pgxpool.Pool, log logger.ILogger) storage.IOrganizeStorage {
	return &organizeRepo{db: db, log: log}
}

func (r *organizeRepo) Create(ctx context.Context, job *models.OrganizeJob) error {
	query := `
		INSERT INTO organize_jobs (id, user_id, input_file_id, new_order, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID, job.NewOrder, job.Status, job.CreatedAt)
	return err
}

func (r *organizeRepo) Update(ctx context.Context, job *models.OrganizeJob) error {
	query := `
		UPDATE organize_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	return err
}

func (r *organizeRepo) GetByID(ctx context.Context, id string) (*models.OrganizeJob, error) {
	query := `
		SELECT id, user_id, input_file_id, new_order, output_file_id, status, created_at
		FROM organize_jobs
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	var job models.OrganizeJob

	err := row.Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.NewOrder,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &job, nil
}
