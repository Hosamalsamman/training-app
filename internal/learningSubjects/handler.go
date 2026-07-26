package learningSubjects

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

func (h *Handler) ListLearningSubjects(c *gin.Context) {

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	learningSubjects, err := h.service.GetAll(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, learningSubjects)
}

func (h *Handler) GetLearningSubject(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	learningSubject, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, learningSubject)
}