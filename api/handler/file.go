package handler

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	// bu path loyihangizga qarab bo'lishi mumkin

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"convertpdfgo/api/models"
	"convertpdfgo/pkg/faylchek"
)

// UploadFile godoc
// @Router       /file/upload [POST]
// @Security     ApiKeyAuth
// @Summary      Upload file
// @Tags         file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Upload file"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UploadFile(c *gin.Context) {
	// Auth optional: user_id bo'lishi shart emas
	var ptrUserID *string
	if uid, exists := c.Get("user_id"); exists {
		if strID, ok := uid.(string); ok && strID != "" {
			ptrUserID = &strID
		}
	}

	// Faylni olish
	fileHeader, err := c.FormFile("file")
	if err != nil {
		handleResponse(c, h.log, "file is required", http.StatusBadRequest, err.Error())
		return
	}

	// Fayl parametrlari
	fileName := fileHeader.Filename
	fileType := filepath.Ext(fileName)
	fileSize := fileHeader.Size
	fileID := uuid.NewString()
	savePath := fmt.Sprintf("uploads/%s%s", fileID, fileType)

	// Fayl kengaytmasini tekshirish
	if faylchek.IsBlacklistedExtension(fileType) {
		handleResponse(c, h.log, "❌ Fayl turi xavfli va yuklash taqiqlangan", http.StatusBadRequest, nil)
		return
	}

	if !faylchek.IsAllowedExtension(fileType) {
		handleResponse(c, h.log, "❌ Bu turdagi fayllarni yuklash ruxsat etilmagan", http.StatusBadRequest, nil)
		return
	}
	// 🔒 Fayl hajmi cheklovi
	const guestMaxSize = 20 * 1024 * 1024      // 30 MB
	const registeredMaxSize = 30 * 1024 * 1024 // 50 MB

	if ptrUserID == nil && fileSize > guestMaxSize {
		handleResponse(c, h.log, "Guests can upload files up to 30MB only", http.StatusBadRequest, nil)
		return
	}
	if ptrUserID != nil && fileSize > registeredMaxSize {
		handleResponse(c, h.log, "Registered users can upload files up to 50MB only", http.StatusBadRequest, nil)
		return
	}

	// Faylni saqlash
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		handleResponse(c, h.log, "failed to save file", http.StatusInternalServerError, err.Error())
		return
	}

	// Bazaga yozish uchun model
	file := models.File{
		ID:         fileID,
		UserID:     ptrUserID, // <-- pointer bo'lishi kerak
		FileName:   fileName,
		FilePath:   savePath,
		FileType:   fileType,
		FileSize:   fileSize,
		UploadedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := h.services.File().Upload(ctx, file)
	if err != nil {
		handleResponse(c, h.log, "upload failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "file uploaded", http.StatusCreated, gin.H{"id": id})
}

// GetFile godoc
// @Router       /file/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get file by ID
// @Tags         file
// @Param        id path string true "File ID"
// @Produce      json
// @Success      200  {object}  models.File
// @Failure      404  {object}  models.Response
func (h Handler) GetFile(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := h.services.File().Get(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "file not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "file found", http.StatusOK, file)
}

// DeleteFile godoc
// @Router       /file/{id} [DELETE]
// @Security     ApiKeyAuth
// @Summary      Delete file by ID
// @Tags         file
// @Param        id path string true "File ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  models.Response
func (h Handler) DeleteFile(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.services.File().Delete(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "delete failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "file deleted", http.StatusOK, gin.H{"id": id})
}

// ListUserFiles godoc
// @Router       /file/list [GET]
// @Security     ApiKeyAuth
// @Summary      List all user's files
// @Tags         file
// @Produce      json
// @Success      200  {array}  models.File
// @Failure      500  {object}  models.Response
func (h Handler) ListUserFiles(c *gin.Context) {
	userID := c.GetString("user_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	files, err := h.services.File().List(ctx, userID)
	if err != nil {
		handleResponse(c, h.log, "list failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "user files", http.StatusOK, files)
}

// CleanupOldFiles godoc
// @Summary      Cleanup old files
// @Description  Admin-only endpoint to delete files older than N days
// @Tags         file
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/files/cleanup [get]
func (h *Handler) CleanupOldFiles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const olderThanDays = 7

	count, err := h.services.File().CleanupOldFiles(ctx, olderThanDays)
	if err != nil {
		handleResponse(c, h.log, "failed to cleanup old files", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "old files cleaned up", http.StatusOK, gin.H{"deleted_files": count})
}
