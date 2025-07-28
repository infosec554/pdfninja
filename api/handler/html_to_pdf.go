package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreateHTMLToPDF godoc
// @Summary      Convert HTML to PDF
// @Description  HTML kontentni PDF formatga o‘tkazish
// @Tags         HTML to PDF
// @Accept       json
// @Produce      json
// @Param        data body models.CreateHTMLToPDFRequest true "HTML to PDF request"
// @Success      201 {object} models.Response{data=string} "job_id qaytadi"
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/html-to-pdf [post]
func (h *Handler) CreateHTMLToPDF(c *gin.Context) {
	var req models.CreateHTMLToPDFRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// HTMLContent bo‘sh bo‘lmasligi kerak
	if req.HTMLContent == "" {
		handleResponse(c, h.log, "html_content is required", http.StatusBadRequest, "empty html_content")
		return
	}

	// Auth bo‘yicha user_id olish
	var userID *string
	if uid := c.GetString("user_id"); uid != "" {
		userID = &uid
	} else {
		userID = nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID, err := h.services.HTMLToPDF().Create(ctx, req, userID)
	if err != nil {
		handleResponse(c, h.log, "conversion failed", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "conversion started", http.StatusCreated, jobID)
}

// GetHTMLToPDFJob godoc
// @Summary      Get HTML→PDF Job Status
// @Description  Konvertatsiya holatini olish
// @Tags         HTML to PDF
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.HTMLToPDFJob
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /api/pdf/html-to-pdf/{id} [get]
func (h *Handler) GetHTMLToPDFJob(c *gin.Context) {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job, err := h.services.HTMLToPDF().GetByID(ctx, jobID)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "job found", http.StatusOK, job)
}
