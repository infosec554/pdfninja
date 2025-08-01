package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateAddPageNumbersJob godoc
// @Router       /api/pdf/add-page-numbers [POST]
// @Summary      Add page numbers to PDF
// @Tags         PDF
// @Accept       json
// @Produce      json
// @Param        data body models.AddPageNumbersRequest true "PDF ID and font settings"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h *Handler) CreateAddPageNumbersJob(c *gin.Context) {
	var req models.AddPageNumbersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// Check if inputFileIDs are provided
	if len(req.InputFileID) == 0 {
		handleResponse(c, h.log, "no input files", http.StatusBadRequest, "input_file_ids required")
		return
	}
	// Handle guest user (if user_id is empty)
	var userID *string
	if uid := c.GetString("user_id"); uid != "" {
		userID = &uid
	} else {
		// For guest user, we pass nil
		userID = nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.AddPageNumber().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create add-page-numbers job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetAddPageNumbersJob godoc
// @Router       /api/pdf/add-page-numbers/{id} [GET]
// @Summary      Get add page numbers job by ID
// @Tags         PDF
// @Param        id path string true "Job ID"
// @Produce      json
// @Success      200 {object} models.AddPageNumberJob
// @Failure      404 {object} models.Response
func (h *Handler) GetAddPageNumbersJob(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.AddPageNumber().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "job found", http.StatusOK, job)
}
