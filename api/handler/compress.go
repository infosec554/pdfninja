package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateCompressJob godoc
// @Router       /api/pdf/compress [POST]
// @Security     ApiKeyAuth
// @Summary      Create compress job
// @Tags         pdf-compress
// @Accept       json
// @Produce      json
// @Param        request body models.CompressRequest true "compress job"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateCompressJob(c *gin.Context) {
	var req models.CompressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Compress().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create compress job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "compress job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetCompressJob godoc
// @Router       /api/pdf/compress/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get compress job by ID
// @Tags         pdf-compress
// @Accept       json
// @Produce      json
// @Param        id path string true "compress job ID"
// @Success      200 {object} models.CompressJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetCompressJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing id", http.StatusBadRequest, "id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Compress().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "failed to get compress job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "compress job fetched", http.StatusOK, job)
}
