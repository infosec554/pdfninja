package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// TranslatePDF godoc
// @Router       /api/pdf/translate [POST]
// @Security     BearerAuth
// @Summary      Translate PDF file to another language
// @Description  Tarjima qilish: PDF fayldagi matnni boshqa tilga o‘girish va yangi PDFga saqlash
// @Tags         pdf-translate
// @Accept       json
// @Produce      json
// @Param        request body models.TranslatePDFRequest true "Translation parameters"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) TranslatePDF(c *gin.Context) {
	var req models.TranslatePDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request payload", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required in token")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobID, err := h.services.TranslatePDF().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create translation job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "translation job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetTranslatePDFJob godoc
// @Router       /api/pdf/translate/{id} [GET]
// @Security     BearerAuth
// @Summary      Get translation job status
// @Description  PDF tarjima qilish jarayonining holatini olish
// @Tags         pdf-translate
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.TranslatePDFJob
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) GetTranslatePDFJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		handleResponse(c, h.log, "job id is required", http.StatusBadRequest, "id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.TranslatePDF().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "translation job status", http.StatusOK, job)
}
