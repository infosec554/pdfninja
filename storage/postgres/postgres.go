package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"test/config"
	"test/pkg/logger"
	"test/storage"
)

type Store struct {
	pool  *pgxpool.Pool
	log   logger.ILogger
	redis storage.IRedisStorage
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger, redis storage.IRedisStorage) (storage.IStorage, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to database", logger.Error(err))
		return nil, err
	}

	absPath, err := filepath.Abs("migrations/postgres")
	if err != nil {
		log.Error("failed to get absolute path for migrations", logger.Error(err))
		return nil, err
	}
	m, err := migrate.New("file://"+absPath, url)
	if err != nil {
		log.Error("migration error", logger.Error(err))
		return nil, err
	}
	if err = m.Up(); err != nil && !strings.Contains(err.Error(), "no change") {
		log.Error("migration up error", logger.Error(err))
		return nil, err
	}

	log.Info("postgres connected and migrated")

	return &Store{
		pool:  pool,
		log:   log,
		redis: redis,
	}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) User() storage.IUserStorage {
	return NewUserRepo(s.pool, s.log)
}

func (s *Store) OTP() storage.IOTPStorage {
	return NewOTPRepo(s.pool, s.log)
}

func (s *Store) Role() storage.IRoleStorage {
	return NewRoleRepo(s.pool, s.log)
}

func (s *Store) Sysuser() storage.ISysuserStorage {
	return NewSysuserRepo(s.pool, s.log)
}

func (s *Store) Redis() storage.IRedisStorage {
	return s.redis
}

func (s *Store) Merge() storage.IMergeStorage {
	return NewMergeRepo(s.pool, s.log)
}

func (s *Store) File() storage.IFileStorage {
	return NewFileRepo(s.pool, s.log)
}

func (s *Store) Split() storage.ISplitStorage {
	return NewSplitRepo(s.pool, s.log)
}

func (s *Store) RemovePage() storage.IRemovePageStorage {
	return NewRemovePageRepo(s.pool, s.log)
}

func (s *Store) ExtractPage() storage.IExtractPageStorage {
	return NewExtractPageRepo(s.pool, s.log)
}

func (s *Store) Organize() storage.IOrganizeStorage {
	return NewOrganizeRepo(s.pool, s.log)
}

func (s *Store) Compress() storage.ICompressStorage {
	return NewCompressRepo(s.pool, s.log)
}

func (s *Store) PDFToJPG() storage.IPDFToJPGStorage {
	return NewPDFToJPGRepo(s.pool, s.log)
}

func (s *Store) Rotate() storage.IRotateStorage {
	return NewRotateRepo(s.pool, s.log)
}

func (s *Store) AddPageNumber() storage.IAddPageNumberStorage {
	return NewAddPageNumberRepo(s.pool, s.log)
}
func (s *Store) Crop() storage.ICropPDFStorage {
	return NewCropRepo(s.pool, s.log)
}

func (s *Store) Unlock() storage.IUnlockPDFStorage {
	return NewUnlockRepo(s.pool, s.log)
}
func (s *Store) Protect() storage.IProtectStorage {
	return NewProtectRepo(s.pool, s.log)
}

func (s *Store) Stat() storage.IStatStorage {
	return NewStatsRepo(s.pool, s.log)
}
func (s *Store) Log() storage.ILogService {
	return NewLogRepo(s.pool, s.log)
}
func (s *Store) JPGToPDF() storage.IJPGToPDFStorage {
	return NewJPGToPDFRepo(s.pool, s.log)
}
func (s *Store) Inspect() storage.IInspectStorage {
	return NewInspectRepo(s.pool, s.log)
}

func (s *Store) TranslatePDF() storage.ITranslatePDFStorage {
	return NewTranslatePDFRepo(s.pool, s.log)
}

func (s *Store) SharedLink() storage.ISharedLinkStorage {
	return NewSharedLinkRepo(s.pool, s.log)
}

func (s *Store) AddHeaderFooter() storage.AddHeaderFooterStorage {
	return NewAddHeaderFooterRepo(s.pool, s.log)
}

func (s *Store) AddBackground() storage.IAddBackgroundStorage {
	return NewAddBackgroundRepo(s.pool, s.log)
}

func (s *Store) DetectBlankPages() storage.IDetectBlankPagesStorage {
	return NewDetectBlankPagesRepo(s.pool, s.log)
}

func (s *Store) QRCode() storage.IQRCodeStorage {
	return NewQRCodeRepo(s.pool, s.log)
}

func (s *Store) PDFTextSearch() storage.IPDFTextSearchStorage {
	return NewPDFTextSearchRepo(s.pool, s.log)
}
