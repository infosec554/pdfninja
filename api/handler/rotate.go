package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateRotateJob godoc
// @Summary      Rotate PDF pages
// @Description  PDF fayl sahifalarini aylantiradi (90, 180 yoki 270 gradus)
// @Tags         PDF Edit
// @Accept       json
// @Produce      json
// @Param        request body models.RotatePDFRequest true "Rotate request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/rotate [post]
// @Security     ApiKeyAuth
func (h *Handler) CreateRotateJob(c *gin.Context) {
	var req models.RotatePDFRequest

	// So‘rovni parse qilish
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// Foydalanuvchi ID sini olish
	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "user_id is missing", http.StatusBadRequest, "auth required")
		return
	}

	// Service chaqiruv
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Rotate().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "rotate job failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "rotate job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetRotateJob godoc
// @Summary      Get Rotate Job info
// @Description  Rotate PDF jarayonining natijasini olish
// @Tags         PDF Edit
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.RotateJob
// @Failure      404 {object} models.Response
// @Router       /api/pdf/rotate/{id} [get]
// @Security     ApiKeyAuth
func (h *Handler) GetRotateJob(c *gin.Context) {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Rotate().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "rotate job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "rotate job found", http.StatusOK, job)
}
