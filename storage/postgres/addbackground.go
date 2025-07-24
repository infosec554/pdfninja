package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type addBackgroundRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewAddBackgroundRepo(db *pgxpool.Pool, log logger.ILogger) storage.IAddBackgroundStorage {
	return &addBackgroundRepo{db: db, log: log}
}

func (r *addBackgroundRepo) Create(ctx context.Context, job *models.AddBackgroundJob) error {
	query := `
		INSERT INTO add_background_jobs
		(id, user_id, input_file_id, background_image_file_id, opacity, position, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID,
		job.UserID,
		job.InputFileID,
		job.BackgroundImageFileID,
		job.Opacity,
		job.Position,
		job.Status,
		job.CreatedAt,
	)
	if err != nil {
		r.log.Error("failed to create add background job", logger.Error(err))
		return err
	}
	return nil
}

func (r *addBackgroundRepo) Update(ctx context.Context, job *models.AddBackgroundJob) error {
	query := `
		UPDATE add_background_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	if err != nil {
		r.log.Error("failed to update add background job", logger.Error(err))
		return err
	}
	return nil
}

func (r *addBackgroundRepo) GetByID(ctx context.Context, id string) (*models.AddBackgroundJob, error) {
	query := `
		SELECT id, user_id, input_file_id, background_image_file_id, opacity, position, output_file_id, status, created_at
		FROM add_background_jobs
		WHERE id = $1
	`
	var job models.AddBackgroundJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.BackgroundImageFileID,
		&job.Opacity,
		&job.Position,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
