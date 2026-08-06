package jobs

import (
	"log"
	"net/http"
	"strconv"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListJobs(c *gin.Context) {

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	jobs, err := h.service.GetAll(clientID)

	if err != nil {

		log.Printf("ListJobs failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetJob(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {

		log.Printf("GetJob failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	job, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *Handler) Create(c *gin.Context) {

	var req models.Job

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات غير صحيحة",
		})
		return
	}

	clientID := 2 // TODO: from JWT

	_,err := h.service.Create(clientID, &req)  // if you need to return the created object use res instead of _ and replace gin.H with res

	if err != nil {

		status, msg := dbutil.TranslateDBError(err)

		c.JSON(status, gin.H{
			"error": msg,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": "تمت الإضافة بنجاح",
	})
}