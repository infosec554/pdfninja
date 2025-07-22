package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateOrganizeJob godoc
// @Router       /api/pdf/organize [POST]
// @Security     ApiKeyAuth
// @Summary      Organize PDF pages (change page order)
// @Tags         pdf-organize
// @Accept       json
// @Produce      json
// @Param        request body models.CreateOrganizeJobRequest true "Organize PDF request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) CreateOrganizeJob(c *gin.Context) {
	var req models.CreateOrganizeJobRequest

	// Request body ni tekshirish
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	// Service chaqiruvi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Organize().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create organize job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "organize job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetOrganizeJob godoc
// @Router       /api/pdf/organize/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get organize job by ID
// @Tags         pdf-organize
// @Accept       json
// @Produce      json
// @Param        id path string true "Organize job ID"
// @Success      200 {object} models.OrganizeJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) GetOrganizeJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing id", http.StatusBadRequest, "id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.Organize().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "failed to get organize job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "organize job fetched", http.StatusOK, job)
}
