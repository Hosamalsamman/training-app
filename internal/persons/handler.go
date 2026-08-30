package persons

import (
	"errors"
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

func (h *Handler) ListPersons(c *gin.Context) {

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	persons, err := h.service.GetAll(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, persons)
}

func (h *Handler) GetPerson(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	// TODO: replace with c.MustGet("clientID").(int)
	clientID := 2

	person, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *Handler) Create(c *gin.Context) {

	var req models.CreatePersonRequest

	// Convert JSON body into our request struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	// client_id was placed into the Gin context
	// by the JWT middleware.
	clientID := c.GetInt("client_id")

	person, err := h.service.Create(clientID, &req)

	if err != nil {
		status, message := dbutil.TranslateDBError(err)

		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// RegisterUser turns an existing person into a user by
// giving him a username, a password and a group.
func (h *Handler) RegisterUser(c *gin.Context) {

	var req models.RegisterUserRequest

	// Convert JSON body into our request struct.
	// The binding tags enforce the required fields
	// and the minimum password length.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	// client_id was placed into the Gin context
	// by the JWT middleware.
	clientID := c.GetInt("client_id")

	person, err := h.service.RegisterUser(clientID, &req)

	if err != nil {

		switch {

		// The person does not exist for this client.
		case errors.Is(err, ErrPersonNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The person already has a username and password.
		case errors.Is(err, ErrAlreadyRegistered):
			c.JSON(http.StatusConflict, gin.H{
				"error": "الشخص مسجل بالفعل كمستخدم.",
			})

		// Anything else is a database error.
		default:
			status, message := dbutil.TranslateDBError(err)

			c.JSON(status, gin.H{
				"error": message,
			})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}
