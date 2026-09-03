package performanceEvaluationCourseParticipantDetails

import (
	"net/http"
	"strconv"

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

func (h *Handler) ListPerformanceEvaluationCourseParticipantDetails(c *gin.Context) {

	clientID := c.GetInt("client_id")

	participantDetails, err := h.service.GetAll(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, participantDetails)
}

func (h *Handler) GetPerformanceEvaluationCourseParticipantDetail(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	clientID := c.GetInt("client_id")

	participantDetail, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, participantDetail)
}
