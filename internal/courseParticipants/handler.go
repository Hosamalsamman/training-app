package courseParticipants

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

func (h *Handler) ListCourseParticipants(c *gin.Context) {

	clientID := c.GetInt("client_id")

	courseParticipants, err := h.service.GetAll(clientID)

	if err != nil {

		log.Printf("ListJobs failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courseParticipants)
}

func (h *Handler) GetCourseParticipant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {

		log.Printf("GetJob failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	clientID := c.GetInt("client_id")

	courseParticipant, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, courseParticipant)
}