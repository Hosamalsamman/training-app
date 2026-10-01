package courseSessionParticipants

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

func (h *Handler) ListCourseSessionParticipants(c *gin.Context) {

	clientID := c.GetInt("client_id")

	courseSessionParticipants, err := h.service.GetAll(clientID)

	if err != nil {

		log.Printf("ListCourseSessionParticipants failed: %v", err)
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

		log.Printf("GetCourseSessionParticipant failed: %v", err)
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

// translateError maps the sentinel errors of the service to
// the HTTP status code and Arabic message each rule deserves.
// Anything else falls through to the DB error translator,
// exactly like the course session handlers do.
func translateError(err error) (int, gin.H) {

	switch {

	// The session referenced by course_session_id is missing,
	// or belongs to another client.
	case errors.Is(err, ErrSessionNotFound):
		return http.StatusNotFound, gin.H{
			"error": "محاضرة المقرر غير موجودة.",
		}

	// The enrollment referenced by course_participant_id is
	// missing, or belongs to another client.
	case errors.Is(err, ErrParticipantNotFound):
		return http.StatusNotFound, gin.H{
			"error": " المتدرب غير موجود.",
		}

	// The attendance row id is missing, or belongs to another
	// client: both are reported the same way so ids cannot be
	// probed across tenants.
	case errors.Is(err, ErrSessionParticipantNotFound):
		return http.StatusNotFound, gin.H{
			"error": "المحاضرة غير موجودة.",
		}

	// The attendance rule: the enrollment belongs to another
	// course than the session.
	case errors.Is(err, ErrParticipantNotInCourse):
		return http.StatusBadRequest, gin.H{
			"error": "لا يمكن تسجيل الحضور: المشارك غير مسجل في هذه الدورة.",
		}

	// Anything else is a database error.
	default:
		status, message := dbutil.TranslateDBError(err)

		return status, gin.H{
			"error": message,
		}
	}
}

// CreateCourseSessionParticipant handles POST
// /new-course-session-participant. Attendance is recorded for
// an enrollment (course participant) and the service
// enforces that the session and the enrollment belong to the
// same course.
func (h *Handler) CreateCourseSessionParticipant(c *gin.Context) {

	var req models.CourseSessionParticipantRequest

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

// UpdateCourseSessionParticipant handles PUT
// /course-session-participant/:id. It is a full update and
// re-applies the same-course rule of creation to the new
// values.
func (h *Handler) UpdateCourseSessionParticipant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req models.CourseSessionParticipantRequest

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

// DeleteCourseSessionParticipant handles DELETE
// /course-session-participant/:id. A row of another client is
// reported exactly like a missing one, so ids cannot be
// probed across tenants.
func (h *Handler) DeleteCourseSessionParticipant(c *gin.Context) {

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
		"message": "تم حذف سجل حضور المشارك بنجاح.",
	})
}