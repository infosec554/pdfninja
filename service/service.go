package service

import (
	"test/pkg/gotenberg"
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

	Compress() CompressService
	PDFToJPG() PDFToJPGService
	Rotate() RotateService
	AddPageNumber() AddPageNumberService
	Crop() CropPDFService
	Unlock() UnlockService
	Protect() ProtectPDFService
	Stat() StatsService
	Log() LogService
	JPGToPDF() JPGToPDFService
	Inspect() InspectService

	SharedLink() SharedLinkService
	AddHeaderFooter() AddHeaderFooterService
	DetectBlank() DetectBlankService
	QRCode() QRCodeService
	PDFToWord() PDFToWordService
	WordToPDF() WordToPDFService
}

type service struct {
	userService            UserService
	otpService             OtpService
	roleService            RoleService
	sysUserService         SysUserService
	mailer                 MailerService
	mergeService           MergeService
	fileService            FileService
	splitService           SplitService
	removepageService      RemovePageService
	extractPageService     ExtractPageService
	compressService        CompressService
	pdfToJPGService        PDFToJPGService
	rotateSrvice           RotateService
	addPageNumberService   AddPageNumberService
	cropPDFService         CropPDFService
	unlockService          UnlockService
	protectPDFService      ProtectPDFService
	statsService           StatsService
	logService             LogService
	jPGToPDF               JPGToPDFService
	inspect                InspectService
	
	sharedLinkService      SharedLinkService
	addHeaderFooterService AddHeaderFooterService
	detectBlankService     DetectBlankService
	qRCodeService          QRCodeService
	pdfToWordService       PDFToWordService
	wordToPDFService       WordToPDFService
}

func New(storage storage.IStorage, log logger.ILogger, mailerCore *mailer.Mailer, redis storage.IRedisStorage, gotClient gotenberg.Client) IServiceManager {
	return &service{
		userService:            NewUserService(storage, log),
		otpService:             NewOtpService(storage, log, mailerCore, redis),
		roleService:            NewRoleService(storage, log),
		sysUserService:         NewSysUserService(storage, log),
		mailer:                 NewMailerService(mailerCore),
		mergeService:           NewMergeService(storage, log),
		fileService:            NewFileService(storage, log),
		splitService:           NewSplitService(storage, log),
		removepageService:      NewRemoveService(storage, log),
		extractPageService:     NewExtractService(storage, log),
		compressService:        NewCompressService(storage, log),
		pdfToJPGService:        NewPDFToJPGService(storage, log),
		rotateSrvice:           NewRotateService(storage, log),
		addPageNumberService:   NewAddPageNumberService(storage, log),
		cropPDFService:         NewCropPDFService(storage, log),
		unlockService:          NewUnlockService(storage, log),
		protectPDFService:      NewProtectPDFService(storage, log),
		statsService:           NewStatsService(storage, log),
		logService:             NewLogService(storage, log),
		jPGToPDF:               NewJPGToPDFService(storage, log),
		inspect:                NewInspectService(storage, log),
	
		sharedLinkService:      NewSharedLinkService(storage, log),
		addHeaderFooterService: NewAddHeaderFooterService(storage, log),
		detectBlankService:     NewDetectBlankService(storage, log),
		qRCodeService:          NewQRCodeService(storage, log),
		pdfToWordService:       NewPDFToWordService(storage, log),
		wordToPDFService:       NewWordToPDFService(storage, log, gotClient),
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

func (s *service) Compress() CompressService {
	return s.compressService
}

func (s *service) PDFToJPG() PDFToJPGService {
	return s.pdfToJPGService
}

func (s *service) Rotate() RotateService {
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

func (s *service) Stat() StatsService {
	return s.statsService
}

func (s *service) Log() LogService {
	return s.logService
}

func (s *service) JPGToPDF() JPGToPDFService {
	return s.jPGToPDF
}

func (s *service) Inspect() InspectService {
	return s.inspect
}


func (s *service) SharedLink() SharedLinkService {
	return s.sharedLinkService
}
func (s *service) AddHeaderFooter() AddHeaderFooterService {
	return s.addHeaderFooterService
}

func (s *service) DetectBlank() DetectBlankService {
	return s.detectBlankService
}

func (s *service) QRCode() QRCodeService {
	return s.qRCodeService
}

func (s *service) PDFToWord() PDFToWordService {
	return s.pdfToWordService
}

func (s *service) WordToPDF() WordToPDFService {
	return s.wordToPDFService
}
