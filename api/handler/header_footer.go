package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// AddHeaderFooter godoc
// @Router       /api/pdf/header-footer [POST]
// @Security     ApiKeyAuth
// @Summary      Add header and/or footer to PDF
// @Description  PDF faylga yuqori yoki quyi matn qo‘shish
// @Tags         pdf-header-footer
// @Accept       json
// @Produce      json
// @Param        request body models.AddHeaderFooterRequest true "Header/Footer parametrlari"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) AddHeaderFooter(c *gin.Context) {
	var req models.CreateAddHeaderFooterRequest
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

	jobID, err := h.services.AddHeaderFooter().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to add header/footer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "header/footer job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetHeaderFooterJob godoc
// @Router       /api/pdf/header-footer/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get header/footer job status
// @Description  PDF faylga header/footer qo‘shish jarayonining holatini olish
// @Tags         pdf-header-footer
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.AddHeaderFooterJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetHeaderFooterJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	job, err := h.services.AddHeaderFooter().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "header/footer job fetched", http.StatusOK, job)
}
