package courseSessionParticipants

import (
	"log"
	"net/http"
	"strconv"
	// "training-app/internal/dbutil"
	// "training-app/internal/models"

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

func (h *Handler) ListCourseSessionParticipants(c *gin.Context) {

	clientID := c.GetInt("client_id")

	courseSessionParticipants, err := h.service.GetAll(clientID)

	if err != nil {

		log.Printf("ListJobs failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courseSessionParticipants)
}

func (h *Handler) GetCourseSessionParticipant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {

		log.Printf("GetJob failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	clientID := c.GetInt("client_id")

	courseSessionParticipant, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, courseSessionParticipant)
}