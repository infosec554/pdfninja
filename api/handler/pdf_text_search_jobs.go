package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
)

// CreatePDFTextSearchJob godoc
// @Router       /api/pdf/text-search [POST]
// @Security     ApiKeyAuth
// @Summary      PDF fayldan matn chiqarib, qidirish uchun tayyorlaydi
// @Description  Foydalanuvchi yuklagan PDF fayldan matnni chiqarib, bazaga saqlaydi
// @Tags         pdf-text-search
// @Accept       json
// @Produce      json
// @Param        request body models.CreatePDFTextSearchRequest true "InputFileID"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreatePDFTextSearchJob(c *gin.Context) {
	var req models.CreatePDFTextSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	var userID *string
	if uid := c.GetString("user_id"); uid != "" {
		userID = &uid
	} else {
		// For guest user, we pass nil
		userID = nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobID, err := h.services.PDFTextSearch().Create(ctx, req.InputFileID, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to create job", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "pdf text search job created", http.StatusCreated, gin.H{"id": jobID})
}

// GetPDFTextSearchJob godoc
// @Router       /api/pdf/text-search/{id} [GET]
// @Security     ApiKeyAuth
// @Summary      PDF matn qidirish ishining holatini ko‘rish
// @Description  Avval yaratilgan PDF matn qidirish ishining holati va chiqarilgan matnni olish
// @Tags         pdf-text-search
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      200 {object} models.PDFTextSearchJob
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
func (h Handler) GetPDFTextSearchJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponse(c, h.log, "missing job id", http.StatusBadRequest, "id is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	job, err := h.services.PDFTextSearch().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, h.log, "job not found", http.StatusNotFound, err.Error())
		return
	}

	handleResponse(c, h.log, "pdf text search job fetched", http.StatusOK, job)
}
