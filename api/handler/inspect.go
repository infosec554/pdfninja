package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateInspectJob godoc
// @Router       /api/pdf/inspect [POST]
// @Security     ApiKeyAuth
// @Summary      Inspect PDF metadata
// @Description  PDF faylning strukturasi, sahifa soni, sarlavha, muallif va boshqa metadata’larni olish
// @Tags         inspect
// @Accept       json
// @Produce      json
// @Param        request body models.InspectRequest true "Inspect request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateInspectJob(c *gin.Context) {
	var req models.InspectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	if req.FileID == "" {
		handleResponse(c, h.log, "missing file_id", http.StatusBadRequest, "file_id is required")
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.Inspect().Create(ctx, req.FileID, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to inspect PDF", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "inspect job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetInspectJob godoc
// @Router       /api/pdf/inspect/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get Inspect Job
// @Description  Inspect ishining natijasini ko‘rish
// @Tags         inspect
// @Produce      json
// @Param        id path string true "Inspect Job ID"
// @Success      200 {object} models.InspectJob
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetInspectJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.Inspect().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "inspect job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "inspect job fetched", http.StatusOK, job)
}
