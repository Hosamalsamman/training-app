package courseParticipants

import (
	"errors"
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

// translateError maps the sentinel errors of the service to
// the HTTP status code and Arabic message each rule deserves.
// Anything else falls through to the DB error translator,
// exactly like the course session handlers do.
func translateError(err error) (int, gin.H) {

	switch {

	// The course referenced by course_id is missing, or
	// belongs to another client.
	case errors.Is(err, ErrCourseNotFound):
		return http.StatusNotFound, gin.H{
			"error": "المقرر غير موجود.",
		}

	// The person referenced by person_id is missing, or
	// belongs to another client.
	case errors.Is(err, ErrPersonNotFound):
		return http.StatusNotFound, gin.H{
			"error": "الموظف غير موجود.",
		}

	// The participant id is missing, or belongs to another
	// client: both are reported the same way so ids cannot
	// be probed across tenants.
	case errors.Is(err, ErrParticipantNotFound):
		return http.StatusNotFound, gin.H{
			"error": "المرشح للدرجة غير موجود.",
		}

	// The candidate rule: the person is not an eligible
	// path grade candidate for the course, or already
	// passed the path grade.
	case errors.Is(err, ErrNotEligiblePathGradeCandidate):
		return http.StatusBadRequest, gin.H{
			"error": "لا يمكن إضافة الموظف كمشارك: يجب أن يكون مرشحاً لنفس درجة المسار (ولنفس الفصل الدراسي في مقررات فصل دراسي) ولم ينجح بها من قبل.",
		}

	// Anything else is a database error.
	default:
		status, message := dbutil.TranslateDBError(err)

		return status, gin.H{
			"error": message,
		}
	}
}

// CreateCourseParticipant handles POST /new-course-participant.
// It enforces the path-grade candidate rule: for a path-grade
// subject (or subject term) course the person must be a path
// grade candidate of the same path grade (and the same term),
// with a final evaluation that is NULL or not one of the
// success ids.
func (h *Handler) CreateCourseParticipant(c *gin.Context) {

	var req models.CourseParticipantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	// client_id was placed into the Gin context
	// by the JWT middleware.
	clientID := c.GetInt("client_id")

	participant, err := h.service.Create(clientID, &req)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, participant)
}

// UpdateCourseParticipant handles PUT /course-participant/:id.
// It is a full update and re-applies the same candidate rule
// of creation to the new values.
func (h *Handler) UpdateCourseParticipant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req models.CourseParticipantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	clientID := c.GetInt("client_id")

	participant, err := h.service.Update(clientID, id, &req)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, participant)
}

// DeleteCourseParticipant handles DELETE
// /course-participant/:id. A participant of another client is
// reported exactly like a missing one, so ids cannot be probed
// across tenants.
func (h *Handler) DeleteCourseParticipant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	clientID := c.GetInt("client_id")

	err = h.service.Delete(clientID, id)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "تم حذف مشارك المقرر بنجاح.",
	})
}