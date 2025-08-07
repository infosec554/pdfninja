package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"convertpdfgo/api/models"
)

// CreateMergeJob godoc
// @Router /api/pdf/merge [post]
// @Summary      Create merge job
// @Tags         pdf-merge
// @Accept       json
// @Produce      json
// @Param        request body models.CreateMergeJobRequest true "merge job"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateMergeJob(c *gin.Context) {
	var req models.CreateMergeJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}
	// Foydalanuvchi ID sini olish (registratsiyasiz foydalanuvchilar uchun ham ishlaydi)
	var userID *string
	if val, ok := c.Get("user_id"); ok {
		if strID, ok := val.(string); ok && strID != "" {
			userID = &strID
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := h.services.Merge().Create(ctx, userID, req.InputFileIDs)
	if err != nil {
		handleResponse(c, h.log, "failed to create merge job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "merge job created", http.StatusCreated, gin.H{"id": id})
}

// GetMergeJob godoc
// @Router /api/pdf/merge/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get merge job
// @Tags         pdf-merge
// @Accept       json
// @Produce      json
// @Param        id path string true "merge job ID"
// @Success      200 {object} models.MergeJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetMergeJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job ID", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Merge().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "merge job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "merge job fetched", http.StatusOK, job)
}

// ProcessMergeJob godoc
// @Router /api/pdf/merge/process/{id} [get]
// @Security ApiKeyAuth
// @Summary Process merge job
// @Tags pdf-merge
// @Param id path string true "merge job ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
func (h Handler) ProcessMergeJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job ID", http.StatusBadRequest, "job ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	outputID, err := h.services.Merge().ProcessJob(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "failed to process merge job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "merge job processed successfully", http.StatusOK, gin.H{
		"output_file_id": outputID,
	})
}
