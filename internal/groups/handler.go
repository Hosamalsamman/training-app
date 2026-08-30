package groups

import (
	"errors"
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

func (h *Handler) ListGroups(c *gin.Context) {

	groups, err := h.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handler) GetGroup(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	group, err := h.service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// ListAllowedGroups returns the groups the authenticated user
// is allowed to assign other users to: every group with an id
// greater than or equal to his own group id.
func (h *Handler) ListAllowedGroups(c *gin.Context) {

	// user_id and client_id were placed into the Gin context
	// by the JWT middleware.
	userID := c.GetInt("user_id")
	clientID := c.GetInt("client_id")

	groups, err := h.service.GetAllowed(clientID, userID)

	if err != nil {

		switch {

		// The authenticated user does not exist for this client.
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The user has no group, so he cannot assign anyone.
		case errors.Is(err, ErrUserHasNoGroup):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "لا يمكن إتمام العملية لأن حسابك غير مرتبط بمجموعة.",
			})

		// Anything else is a database error.
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, groups)
}
