package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateCropJob godoc
// @Summary      Crop PDF pages
// @Description  PDF sahifalarini belgilangan tomonlardan qirqadi (top, bottom, left, right)
// @Security     ApiKeyAuth
// @Tags         PDF Edit
// @Accept       json
// @Produce      json
// @Param        request body models.CropPDFRequest true "Crop request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/crop [post]
func (h *Handler) CreateCropJob(c *gin.Context) {
	var req models.CropPDFRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "user_id is missing", http.StatusBadRequest, "auth required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Crop().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "crop job failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "crop job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetCropJob godoc
// @Summary      Get Crop Job info
// @Description  Crop PDF jarayonining natijasini olish
// @Tags         PDF Edit
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.CropPDFJob
// @Failure      404 {object} models.Response
// @Router       /api/pdf/crop/{id} [get]
// @Security     ApiKeyAuth
func (h *Handler) GetCropJob(c *gin.Context) {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Crop().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "crop job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "crop job found", http.StatusOK, job)
}
