package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateJpgToPdf godoc
// @Router       /api/pdf/jpg-to-pdf [POST]
// @Security     ApiKeyAuth
// @Summary      Convert multiple JPG files to single PDF
// @Description  Bitta PDF faylga bir nechta JPG rasmlarni birlashtiradi
// @Tags         pdf-jpg-to-pdf
// @Accept       json
// @Produce      json
// @Param        request body models.CreateJpgToPdfRequest true "List of image file IDs"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateJpgToPdf(c *gin.Context) {
	var req models.CreateJpgToPdfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	jobID, err := h.services.JpgToPdf().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to convert JPG to PDF", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "jpg to pdf job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetJpgToPdfJob godoc
// @Router       /api/pdf/jpg-to-pdf/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get JPG to PDF job status
// @Description  Yaratilgan JPG to PDF job holatini olish
// @Tags         pdf-jpg-to-pdf
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.JpgToPdfJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetJpgToPdfJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.JpgToPdf().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "jpg to pdf job fetched", http.StatusOK, job)
}
