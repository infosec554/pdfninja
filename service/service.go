package service

import (
	"test/pkg/logger"
	"test/pkg/mailer"
	"test/storage"
)

type IServiceManager interface {
	User() UserService
	Otp() OtpService
	Role() RoleService
	SysUser() SysUserService
	Mailer() MailerService

	File() FileService
	Merge() MergeService
	Split() SplitService // ✅ QO‘SH!
	RemovePage() RemovePageService
	ExtractPage() ExtractPageService
	Organize() OrganizeService
	Compress() CompressService
	JpgToPdf() JpgToPdfService
	PDFToJPG() PDFToJPGService
	PdfToWord() PdfToWordService
	Rotate() RoatateService
	AddPageNumber() AddPageNumberService
	Crop() CropPDFService
	Unlock() UnlockService
	Protect() ProtectPDFService
}

type service struct {
	userService          UserService
	otpService           OtpService
	roleService          RoleService
	sysUserService       SysUserService
	mailer               MailerService
	mergeService         MergeService
	fileService          FileService
	splitService         SplitService
	removepageService    RemovePageService
	extractPageService   ExtractPageService
	organizeService      OrganizeService
	compressService      CompressService
	jpgToPdfService      JpgToPdfService
	pdfToJPGService      PDFToJPGService
	pdfToWordService     PdfToWordService
	rotateSrvice         RoatateService
	addPageNumberService AddPageNumberService
	cropPDFService       CropPDFService
	unlockService        UnlockService
	protectPDFService    ProtectPDFService
}

func New(storage storage.IStorage, log logger.ILogger, mailerCore *mailer.Mailer, redis storage.IRedisStorage) IServiceManager {
	return &service{
		userService:          NewUserService(storage, log),
		otpService:           NewOtpService(storage, log, mailerCore, redis),
		roleService:          NewRoleService(storage, log),
		sysUserService:       NewSysUserService(storage, log),
		mailer:               NewMailerService(mailerCore),
		mergeService:         NewMergeService(storage, log),
		fileService:          NewFileService(storage, log),
		splitService:         NewSplitService(storage, log),
		removepageService:    NewRemoveService(storage, log),
		extractPageService:   NewExtractService(storage, log),
		organizeService:      NewOrganizeService(storage, log),
		compressService:      NewCompressService(storage, log),
		jpgToPdfService:      NewJpgToPdfService(storage, log),
		pdfToJPGService:      NewPDFToJPGService(storage, log),
		pdfToWordService:     NewPdfToWordService(storage, log),
		rotateSrvice:         NewRotateService(storage, log),
		addPageNumberService: NewAddPageNumberService(storage, log),
		cropPDFService:       NewCropPDFService(storage, log),
		unlockService:        NewUnlockService(storage, log),
		protectPDFService:    NewProtectPDFService(storage, log),
	}
}

func (s *service) User() UserService {
	return s.userService
}

func (s *service) Otp() OtpService {
	return s.otpService
}

func (s *service) Role() RoleService {
	return s.roleService
}

func (s *service) SysUser() SysUserService {
	return s.sysUserService
}

func (s *service) Mailer() MailerService {
	return s.mailer
}

func (s *service) Merge() MergeService {
	return s.mergeService
}

func (s *service) File() FileService {
	return s.fileService
}

func (s *service) Split() SplitService {
	return s.splitService
}

func (s *service) RemovePage() RemovePageService {
	return s.removepageService
}

func (s *service) ExtractPage() ExtractPageService {
	return s.extractPageService
}

func (s *service) Organize() OrganizeService {
	return s.organizeService
}
func (s *service) Compress() CompressService {
	return s.compressService
}

func (s *service) JpgToPdf() JpgToPdfService {
	return s.jpgToPdfService
}

func (s *service) PDFToJPG() PDFToJPGService {
	return s.pdfToJPGService
}

func (s *service) PdfToWord() PdfToWordService {
	return s.pdfToWordService
}

func (s *service) Rotate() RoatateService {
	return s.rotateSrvice
}

func (s *service) AddPageNumber() AddPageNumberService {
	return s.addPageNumberService
}

func (s *service) Crop() CropPDFService {
	return s.cropPDFService
}

func (s *service) Unlock() UnlockService {
	return s.unlockService
}

func (s *service) Protect() ProtectPDFService {
	return s.protectPDFService
}
