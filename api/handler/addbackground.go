package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateAddBackground godoc
// @Router       /api/pdf/add-background [POST]
// @Security     ApiKeyAuth
// @Summary      Add background image to PDF
// @Description  PDF faylga orqa fon rasmi qo‘shadi
// @Tags         pdf-add-background
// @Accept       json
// @Produce      json
// @Param        request body models.CreateAddBackgroundRequest true "Add background image request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateAddBackground(c *gin.Context) {
	var req models.CreateAddBackgroundRequest
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

	jobID, err := h.services.AddBackground().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to add background image", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "background image added", http.StatusCreated, gin.H{"id": jobID})
}

// GetAddBackgroundJob godoc
// @Router       /api/pdf/add-background/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get Add Background job status
// @Description  Qo‘shilgan orqa fon rasmi job holatini olish
// @Tags         pdf-add-background
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.AddBackgroundJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetAddBackgroundJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.AddBackground().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "add background job fetched", http.StatusOK, job)
}
