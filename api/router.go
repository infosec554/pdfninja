package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "convertpdfgo/api/docs"
	"convertpdfgo/api/handler"
	"convertpdfgo/pkg/logger"
	"convertpdfgo/pkg/middleware"
	"convertpdfgo/service"
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

	r.POST("/signup", h.SignUp)
	r.POST("/login", h.Login)
	r.POST("/change-password", h.AuthorizerMiddleware, h.ChangePassword)
	r.POST("/auth/google", h.GoogleAuth)
	r.POST("/auth/github", h.GithubAuth)
	r.POST("/auth/facebook", h.GithubAuth)

	r.GET("/me", h.AuthorizerMiddleware, h.GetMyProfile)

	// === Logs (adminlar uchun) ===
	admin := r.Group("/admin")
	admin.Use(h.AuthorizerMiddleware, h.AdminMiddleware)
	{
		admin.GET("/logs/:id", h.GetLogsByJobID)
	}

	stats := r.Group("/stats")
	stats.Use(h.AuthorizerMiddleware)
	{
		stats.GET("/user", h.GetUserStats)
	}

	// === Fayllar (token kerak, chunki user_id kerak) ===
	r.POST("/file/upload", h.UploadFile)

	file := r.Group("/file")

	file.Use(h.AuthorizerMiddleware)

	{
		file.GET("/:id", h.GetFile)
		file.DELETE("/:id", h.DeleteFile)
		file.GET("/list", h.ListUserFiles)
		file.GET("/cleanup", h.AdminMiddleware, h.CleanupOldFiles)
	}

	// === PDF xizmatlari (token shart emas — optional auth) ===
	pdf := r.Group("/api/pdf")

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

		pdf.POST("/compress", h.CreateCompressJob)
		pdf.GET("/compress/:id", h.GetCompressJob)

		pdf.POST("/jpg-to-pdf", h.CreateJPGToPDF)
		pdf.GET("/jpg-to-pdf/:id", h.GetJPGToPDFJob)

		pdf.POST("/pdf-to-jpg", h.CreatePDFToJPG)
		pdf.GET("/pdf-to-jpg/:id", h.GetPDFToJPG)

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

		pdf.POST("/share", h.CreateSharedLink)
		pdf.GET("/share/:token", h.GetSharedLink)

		pdf.POST("/pdf-to-word", h.CreatePDFToWordJob)
		pdf.GET("/pdf-to-word/:id", h.GetPDFToWordJob)

		pdf.POST("/word-to-pdf", h.CreateWordToPDF)
		pdf.GET("/word-to-pdf/:id", h.GetWordToPDFJob)

		pdf.POST("/excel-to-pdf", h.CreateExcelToPDF)
		pdf.GET("/excel-to-pdf/:id", h.GetExcelToPDFJob)

		pdf.POST("/ppt-to-pdf", h.CreatePowerPointToPDF)
		pdf.GET("/ppt-to-pdf/:id", h.GetPowerPointToPDFJob)

		pdf.POST("/watermark", h.AddTextWatermark)
		pdf.GET("/watermark/:id", h.GetWatermarkJob)

	}

	return r
}
