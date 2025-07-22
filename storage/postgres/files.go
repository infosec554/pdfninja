package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type fileRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func  NewFileRepo(db *pgxpool.Pool, log logger.ILogger) storage.IFileStorage {
	return &fileRepo{
		db:  db,
		log: log,
	}
}

// Save - faylni DBga yozish
func (f *fileRepo) Save(ctx context.Context, file models.File) (string, error) {
	query := `
		INSERT INTO files (id, user_id, file_name, file_path, file_type, file_size, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := f.db.Exec(ctx, query,
		file.ID, file.UserID, file.FileName, file.FilePath,
		file.FileType, file.FileSize, file.UploadedAt)
	if err != nil {
		f.log.Error("DB insert error", logger.Error(err))
		return "", err
	}
	return file.ID, nil
}

// GetByID - faylni ID bo‘yicha olish
func (f *fileRepo) GetByID(ctx context.Context, id string) (models.File, error) {
	var file models.File
	query := `
		SELECT id, user_id, file_name, file_path, file_type, file_size, uploaded_at
		FROM files WHERE id = $1
	`
	err := f.db.QueryRow(ctx, query).Scan(
		&file.ID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType, &file.FileSize, &file.UploadedAt)
	if err != nil {
		f.log.Error("failed to fetch file", logger.Error(err))
		return models.File{}, err
	}
	return file, nil
}

// Delete - faylni ID bo‘yicha o‘chirish
func (f *fileRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := f.db.Exec(ctx, query, id)
	if err != nil {
		f.log.Error("failed to delete file", logger.Error(err))
	}
	return err
}

// ListByUser - foydalanuvchiga tegishli fayllar
func (f *fileRepo) ListByUser(ctx context.Context, userID string) ([]models.File, error) {
	query := `
		SELECT id, user_id, file_name, file_path, file_type, file_size, uploaded_at
		FROM files WHERE user_id = $1 ORDER BY uploaded_at DESC
	`
	rows, err := f.db.Query(ctx, query, userID)
	if err != nil {
		f.log.Error("failed to fetch user files", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		err = rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType, &file.FileSize, &file.UploadedAt)
		if err != nil {
			f.log.Error("error scanning row", logger.Error(err))
			continue
		}
		files = append(files, file)
	}
	return files, nil
}
