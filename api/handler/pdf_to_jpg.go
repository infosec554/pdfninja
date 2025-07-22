package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreatePDFToJPG godoc
// @Summary      Convert PDF to JPG
// @Description  Bir PDF faylni sahifalarini JPG formatga o‘girish
// @Tags         pdf-to-jpg
// @Accept       json
// @Produce      json
// @Param        request body models.PDFToJPGRequest true "PDF file ID"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/pdf/pdf-to-jpg [post]
// @Security     ApiKeyAuth
func (h *Handler) CreatePDFToJPG(c *gin.Context) {
	var req models.PDFToJPGRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobID, err := h.services.PDFToJPG().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "conversion failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "PDF to JPG conversion started", http.StatusCreated, gin.H{"job_id": jobID})
}

// GetPDFToJPG godoc
// @Summary      Get PDF to JPG conversion job
// @Description  Job ID bo‘yicha JPG sahifalar holatini olish
// @Tags         pdf-to-jpg
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200  {object}  models.PDFToJPGJob
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/pdf/pdf-to-jpg/{id} [get]
// @Security     ApiKeyAuth
func (h *Handler) GetPDFToJPG(c *gin.Context) {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.PDFToJPG().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "job fetched", http.StatusOK, job)
}
