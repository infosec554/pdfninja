package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "test/api/docs"
	"test/api/handler"
	"test/pkg/logger"
	"test/service"
	"test/pkg/middleware"
)

// @title           Auth API
// @version         1.0
// @description     Authentication and Authorization API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceManager, log logger.ILogger) *gin.Engine {
	h := handler.New(services, log)
	r := gin.Default()

	// === Swagger ===
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.RateLimiterMiddleware())

	// === OTP ===
	r.POST("/otp/send", h.SendOTP)
	r.POST("/otp/confirm", h.ConfirmOTP)

	// === Auth ===
	r.POST("/signup", h.SignUp)
	r.POST("/login", h.Login)
	r.GET("/me", h.AuthorizerMiddleware, h.GetMyProfile)

	// === Logs (faqat adminlar uchun yoki auth foydalanuvchilar) ===
	admin := r.Group("/admin")
	admin.Use(h.AuthorizerMiddleware, h.AdminMiddleware)
	{
		admin.GET("/logs/:id", h.GetLogsByJobID)
	}

	// === Role (faqat adminlar uchun) ===
	role := r.Group("/role")
	role.Use(h.AuthorizerMiddleware, h.AdminMiddleware)
	{
		role.POST("/", h.CreateRole)
		role.PUT("/:id", h.UpdateRole)
		role.GET("/", h.ListRoles)
	}

	stats := r.Group("/stats")
	stats.Use(h.AuthorizerMiddleware)
	{
		stats.GET("/user", h.GetUserStats)
	}

	// === SysUser (admin uchun) ===
	sysuser := r.Group("/sysuser")
	sysuser.Use(h.AuthorizerMiddleware, h.AdminMiddleware)
	{
		sysuser.POST("/", h.CreateSysUser)
	}

	// === Fayllar (user_id ni token orqali oladi) ===
	// === Fayllar (user_id ni token orqali oladi) ===
	file := r.Group("/file")
	file.Use(h.AuthorizerMiddleware)
	{
		file.POST("/upload", h.UploadFile)
		file.GET("/:id", h.GetFile)
		file.DELETE("/:id", h.DeleteFile)
		file.GET("/list", h.ListUserFiles)

		// cleanup faqat adminlar uchun
		file.GET("/cleanup", h.AdminMiddleware, h.CleanupOldFiles)
	}

	// === PDF xizmatlari ===
	pdf := r.Group("/api/pdf")
	pdf.Use(h.AuthorizerMiddleware)
	{
		pdf.POST("/merge", h.CreateMergeJob)
		pdf.GET("/merge/:id", h.GetMergeJob)
		pdf.GET("/merge/process/:id", h.ProcessMergeJob)

		pdf.POST("/split", h.CreateSplitJob)
		pdf.GET("/split/:id", h.GetSplitJob)

		pdf.POST("/removepage", h.CreateRemovePagesJob)
		pdf.GET("/removepage/:id", h.GetRemovePagesJob)

		pdf.POST("/extract", h.CreateExtractJob)
		pdf.GET("/extract/:id", h.GetExtractJob)

		pdf.POST("/organize", h.CreateOrganizeJob)
		pdf.GET("/organize/:id", h.GetOrganizeJob)

		pdf.POST("/compress", h.CreateCompressJob)
		pdf.GET("/compress/:id", h.GetCompressJob)

		pdf.POST("/jpg-to-pdf", h.CreateJpgToPdf)
		pdf.GET("/jpg-to-pdf/:id", h.GetJpgToPdfJob)

		pdf.POST("/pdf-to-jpg", h.CreatePDFToJPG)
		pdf.GET("/pdf-to-jpg/:id", h.GetPDFToJPG)

		pdf.POST("/pdf-to-word", h.CreatePdfToWordJob)
		pdf.GET("/pdf-to-word/:id", h.GetPdfToWordJob)

		pdf.POST("/rotate", h.CreateRotateJob)
		pdf.GET("/rotate/:id", h.GetRotateJob)

		pdf.POST("/crop", h.CreateCropJob)
		pdf.GET("/crop/:id", h.GetCropJob)

		pdf.POST("/unlock", h.CreateUnlockJob)
		pdf.GET("/unlock/:id", h.GetUnlockJob)

		pdf.POST("/protect", h.CreateProtectJob)
		pdf.GET("/protect/:id", h.GetProtectJob)

		pdf.POST("/add-page-numbers", h.CreateAddPageNumbersJob)
		pdf.GET("/add-page-numbers/:id", h.GetAddPageNumbersJob)
	}

	return r
}
