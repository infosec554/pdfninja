package storage

import (
	"context"
	"time"

	"test/api/models"
)

type IStorage interface {
	User() IUserStorage
	OTP() IOTPStorage
	Role() IRoleStorage
	Sysuser() ISysuserStorage
	Redis() IRedisStorage
	Close()
	Merge() IMergeStorage
	File() IFileStorage
	Split() ISplitStorage
	RemovePage() IRemovePageStorage
	ExtractPage() IExtractPageStorage

	Compress() ICompressStorage
	PDFToJPG() IPDFToJPGStorage
	Rotate() IRotateStorage
	AddPageNumber() IAddPageNumberStorage
	Crop() ICropPDFStorage
	Unlock() IUnlockPDFStorage
	Protect() IProtectStorage
	Stat() IStatStorage
	Log() ILogService
	JPGToPDF() IJPGToPDFStorage
	Inspect() IInspectStorage

	SharedLink() ISharedLinkStorage
	AddHeaderFooter() AddHeaderFooterStorage
	DetectBlankPages() IDetectBlankPagesStorage
	QRCode() IQRCodeStorage

	PDFToWord() IPDFToWordStorage
	WordToPDF() IWordToPDFStorage
}

type IUserStorage interface {
	Create(ctx context.Context, req models.CreateUser) (string, error)
	GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type IOTPStorage interface {
	Create(ctx context.Context, email string, code string, expiresAt time.Time) (string, error)
	GetUnconfirmedByID(ctx context.Context, id string) (email string, code string, expiresAt time.Time, err error)
	UpdateStatusToConfirmed(ctx context.Context, id string) error
	GetByIDAndEmail(ctx context.Context, id string, email string) (bool, error)
}

type IRoleStorage interface {
	Create(ctx context.Context, name string, createdBy string) (string, error)
	Update(ctx context.Context, id, name string) error
	GetAll(ctx context.Context) ([]models.Role, error)
	Exists(ctx context.Context, id string) (bool, error)
}

type ISysuserStorage interface {
	GetByPhone(ctx context.Context, phone string) (id, hashedPassword string, status string, err error)
	Create(ctx context.Context, name, phone, hashedPassword, createdBy string) (string, error)
	AssignRoles(ctx context.Context, sysuserID string, roleIDs []string) error
}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type IFileStorage interface {
	Save(ctx context.Context, file models.File) (string, error)
	GetByID(ctx context.Context, id string) (models.File, error)
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]models.File, error)

	GetOldFiles(ctx context.Context, olderThanDays int) ([]models.OldFile, error)
	DeleteByID(ctx context.Context, id string) error
}
type IMergeStorage interface {
	Create(ctx context.Context, job *models.MergeJob) error
	GetByID(ctx context.Context, id string) (*models.MergeJob, error)
	AddInputFiles(ctx context.Context, jobID string, fileIDs []string) error
	GetInputFiles(ctx context.Context, jobID string) ([]string, error)
	Update(ctx context.Context, job *models.MergeJob) error // 👉 YANGI QATOR

}

type ISplitStorage interface {
	Create(ctx context.Context, job *models.SplitJob) error
	GetByID(ctx context.Context, id string) (*models.SplitJob, error)
	UpdateOutputFiles(ctx context.Context, jobID string, outputIDs []string) error
}

type IRemovePageStorage interface {
	Create(ctx context.Context, job *models.RemoveJob) error
	Update(ctx context.Context, job *models.RemoveJob) error
	GetByID(ctx context.Context, id string) (*models.RemoveJob, error)
}

type IExtractPageStorage interface {
	Create(ctx context.Context, job *models.ExtractJob) error
	Update(ctx context.Context, job *models.ExtractJob) error
	GetByID(ctx context.Context, id string) (*models.ExtractJob, error)
}

type IOrganizeStorage interface {
	Create(ctx context.Context, job *models.OrganizeJob) error
	Update(ctx context.Context, job *models.OrganizeJob) error
	GetByID(ctx context.Context, id string) (*models.OrganizeJob, error)
}

type ICompressStorage interface {
	Create(ctx context.Context, job *models.CompressJob) error
	Update(ctx context.Context, job *models.CompressJob) error
	GetByID(ctx context.Context, id string) (*models.CompressJob, error)
}

type IPDFToJPGStorage interface {
	Create(ctx context.Context, job *models.PDFToJPGJob) error
	Update(ctx context.Context, job *models.PDFToJPGJob) error
	GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error)
}

type IRotateStorage interface {
	Create(ctx context.Context, job *models.RotateJob) error
	GetByID(ctx context.Context, id string) (*models.RotateJob, error)
	Update(ctx context.Context, job *models.RotateJob) error
}

type IAddPageNumberStorage interface {
	Create(ctx context.Context, job *models.AddPageNumberJob) error
	GetByID(ctx context.Context, id string) (*models.AddPageNumberJob, error)
	Update(ctx context.Context, job *models.AddPageNumberJob) error
}

type ICropPDFStorage interface {
	Create(ctx context.Context, job *models.CropPDFJob) error
	GetByID(ctx context.Context, id string) (*models.CropPDFJob, error)
	Update(ctx context.Context, job *models.CropPDFJob) error
}

type IUnlockPDFStorage interface {
	Create(ctx context.Context, job *models.UnlockPDFJob) error
	GetByID(ctx context.Context, id string) (*models.UnlockPDFJob, error)
	Update(ctx context.Context, job *models.UnlockPDFJob) error
}
type IProtectStorage interface {
	Create(ctx context.Context, job *models.ProtectPDFJob) error
	GetByID(ctx context.Context, id string) (*models.ProtectPDFJob, error)
	Update(ctx context.Context, job *models.ProtectPDFJob) error
}

type IStatStorage interface {
	GetUserStats(ctx context.Context, userID string) (models.UserStats, error)
}

type ILogService interface {
	GetLogsByJobID(ctx context.Context, jobID string) ([]models.Log, error)
}
type IJPGToPDFStorage interface {
	Create(ctx context.Context, job *models.JPGToPDFJob) error
	GetByID(ctx context.Context, id string) (*models.JPGToPDFJob, error)
	UpdateStatusAndOutput(ctx context.Context, id, status, outputFileID string) error
}

type IInspectStorage interface {
	Create(ctx context.Context, job *models.InspectJob) error
	GetByID(ctx context.Context, id string) (*models.InspectJob, error)
}


type ISharedLinkStorage interface {
	Create(ctx context.Context, req *models.SharedLink) error
	GetByToken(ctx context.Context, token string) (*models.SharedLink, error)
}

type AddHeaderFooterStorage interface {
	Create(ctx context.Context, job *models.AddHeaderFooterJob) error
	Update(ctx context.Context, job *models.AddHeaderFooterJob) error
	GetByID(ctx context.Context, id string) (*models.AddHeaderFooterJob, error)
}

type IDetectBlankPagesStorage interface {
	Create(ctx context.Context, job *models.DetectBlankPagesJob) error
	Update(ctx context.Context, job *models.DetectBlankPagesJob) error
	GetByID(ctx context.Context, id string) (*models.DetectBlankPagesJob, error)
}

type IQRCodeStorage interface {
	Create(ctx context.Context, job *models.QRCodeJob) error
	Update(ctx context.Context, job *models.QRCodeJob) error
	GetByID(ctx context.Context, id string) (*models.QRCodeJob, error)
}



type IPDFToWordStorage interface {
	Create(ctx context.Context, job *models.PDFToWordJob) error
	GetByID(ctx context.Context, id string) (*models.PDFToWordJob, error)
	Update(ctx context.Context, job *models.PDFToWordJob) error
}
type IWordToPDFStorage interface {
	Create(ctx context.Context, job *models.WordToPDFJob) error
	GetByID(ctx context.Context, id string) (*models.WordToPDFJob, error)
	Update(ctx context.Context, job *models.WordToPDFJob) error
}
