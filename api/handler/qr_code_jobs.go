package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateQRCodeJob godoc
// @Router       /api/pdf/qr-code [POST]
// @Security     ApiKeyAuth
// @Summary      Add QR code to PDF
// @Description  PDF faylga QR kod joylashtirish
// @Tags         pdf-qr-code
// @Accept       json
// @Produce      json
// @Param        request body models.CreateQRCodeRequest true "QR Code request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateQRCodeJob(c *gin.Context) {
	var req models.CreateQRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, "user_id required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	jobID, err := h.services.QRCode().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create qr code job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "qr code job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetQRCodeJob godoc
// @Router       /api/pdf/qr-code/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      Get QR code job status
// @Description  QR kod joylashtirish ishining holatini olish
// @Tags         pdf-qr-code
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.QRCodeJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetQRCodeJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.QRCode().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "qr code job fetched", http.StatusOK, job)
}
