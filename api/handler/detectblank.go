package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateDetectBlankPagesJob godoc
// @Router       /api/pdf/detect-blank [POST]
// @Security     ApiKeyAuth
// @Summary      Detect blank pages in PDF
// @Description  PDF fayldagi bo‘sh sahifalarni aniqlaydi
// @Tags         pdf-detect-blank
// @Accept       json
// @Produce      json
// @Param        inputFileID body models.DetectBlankPagesRequest true "Input PDF file ID"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateDetectBlankPagesJob(c *gin.Context) {
	var req models.DetectBlankPagesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobID, err := h.services.DetectBlank().Create(ctx, req.InputFileID, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create detect blank pages job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "detect blank pages job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetDetectBlankPagesJob godoc
// @Router       /api/pdf/detect-blank/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get Detect Blank Pages Job Status
// @Description  Detect blank pages job holatini oladi
// @Tags         pdf-detect-blank
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.DetectBlankPagesJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetDetectBlankPagesJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.DetectBlank().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "detect blank pages job fetched", http.StatusOK, job)
}
