package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreatePdfToWordJob godoc
// @Router       /api/pdf/pdf-to-word [POST]
// @Security     ApiKeyAuth
// @Summary      Convert PDF to Word
// @Description  Yuklangan PDF faylni DOCX formatga o‘zgartiradi
// @Tags         PDF
// @Accept       json
// @Produce      json
// @Param        request body models.PDFToWordRequest true "PDF to Word conversion data"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) CreatePdfToWordJob(c *gin.Context) {
	var req models.PDFToWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "user_id is missing", http.StatusBadRequest, "unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.PdfToWord().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "pdf to word job created", http.StatusCreated, gin.H{"job_id": jobID})
}

// GetPdfToWordJob godoc
// @Router       /api/pdf/pdf-to-word/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get PDF to Word job by ID
// @Description  Konvertatsiya jarayonining natijasini olish
// @Tags         PDF
// @Param        id path string true "Job ID"
// @Produce      json
// @Success      200 {object} models.PDFToWordJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) GetPdfToWordJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "job id is required", http.StatusBadRequest, "missing job id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.PdfToWord().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "job found", http.StatusOK, job)
}
