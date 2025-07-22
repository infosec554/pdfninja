package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateSplitJob godoc
// @Router       /api/pdf/split [post]
// @Security     ApiKeyAuth
// @Summary      Create split job
// @Tags         pdf-split
// @Accept       json
// @Produce      json
// @Param        request body models.CreateSplitJobRequest true "split job"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateSplitJob(c *gin.Context) {
	var req models.CreateSplitJobRequest
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

	jobID, err := h.services.Split().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create split job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "split job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetSplitJob godoc
// @Router       /pdf/split/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get split job by ID
// @Tags         pdf-split
// @Accept       json
// @Produce      json
// @Param        id path string true "split job ID"
// @Success      200 {object} models.SplitJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetSplitJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing id", http.StatusBadRequest, "id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Split().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "failed to get split job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "split job fetched", http.StatusOK, job)
}
