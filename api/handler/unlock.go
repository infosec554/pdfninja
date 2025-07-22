package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateUnlockJob godoc
// @Summary      Unlock PDF file
// @Description  Qulflangan PDF faylni ochish (parolsiz qilish)
// @Tags         PDF Unlock
// @Accept       json
// @Produce      json
// @Param        request body models.UnlockPDFRequest true "Unlock PDF request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/unlock [post]
// @Security     ApiKeyAuth
func (h *Handler) CreateUnlockJob(c *gin.Context) {
	var req models.UnlockPDFRequest

	// So‘rovni tekshirish
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// Auth orqali user_id olish
	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "user_id is required", http.StatusBadRequest, "auth required")
		return
	}

	// Kontekst yaratish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Unlock().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "unlock job failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "unlock job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetUnlockJob godoc
// @Summary      Get Unlock Job status
// @Description  Unlock PDF ish holatini ko‘rish
// @Tags         PDF Unlock
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.UnlockPDFJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/unlock/{id} [get]
// @Security     ApiKeyAuth
func (h *Handler) GetUnlockJob(c *gin.Context) {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Unlock().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "unlock job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "unlock job found", http.StatusOK, job)
}
