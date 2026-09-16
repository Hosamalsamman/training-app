package courseSessions

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

func (h *Handler) ListCourseSessions(c *gin.Context) {

	clientID := c.GetInt("client_id")

	courseSessions, err := h.service.GetAll(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courseSessions)
}

func (h *Handler) GetCourseSession(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	clientID := c.GetInt("client_id")

	courseSession, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, courseSession)
}

// translateError maps the sentinel errors of the service to
// the HTTP status code and Arabic message each rule deserves.
// Anything else falls through to the DB error translator,
// exactly like the course create flow does.
func translateError(err error) (int, gin.H) {

	switch {

	// The time-order rule: start_time before end_time.
	case errors.Is(err, ErrInvalidSessionTimes):
		return http.StatusBadRequest, gin.H{
			"error": "وقت بدء المحاضرة يجب أن يكون قبل وقت انتهائها.",
		}

	// The date-range rule: session_date within the course
	// starting and ending dates.
	case errors.Is(err, ErrSessionOutsideCourseDates):
		return http.StatusBadRequest, gin.H{
			"error": "يجب أن يكون تاريخ المحاضرة ضمن تاريخي بدء وانتهاء المقرر.",
		}

	// The course referenced by course_id exists but is not
	// executed: sessions can be added to executed courses
	// only.
	case errors.Is(err, ErrCourseNotExecuted):
		return http.StatusBadRequest, gin.H{
			"error": "لا يمكن إضافة محاضرات إلا للمقررات المنفذة.",
		}

	// The course referenced by course_id is missing, or
	// belongs to another client.
	case errors.Is(err, ErrCourseNotFound):
		return http.StatusNotFound, gin.H{
			"error": "المقرر غير موجود.",
		}

	// The session id is missing, or belongs to another
	// client.
	case errors.Is(err, ErrSessionNotFound):
		return http.StatusNotFound, gin.H{
			"error": "محاضرة المقرر غير موجودة.",
		}

	// Anything else is a database error.
	default:
		status, message := dbutil.TranslateDBError(err)

		return status, gin.H{
			"error": message,
		}
	}
}

// CreateCourseSession handles POST /new-course-session. It
// creates a session and enforces that its date falls within
// the starting and ending dates of the course it belongs to,
// and that start_time is before end_time.
func (h *Handler) CreateCourseSession(c *gin.Context) {

	var req models.CourseSessionRequest

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

	session, err := h.service.Create(clientID, &req)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, session)
}

// UpdateCourseSession handles PUT /course-session/:id. It is a
// full update and re-applies the same course-date and
// time-order rules of creation to the new values.
func (h *Handler) UpdateCourseSession(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req models.CourseSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	clientID := c.GetInt("client_id")

	session, err := h.service.Update(clientID, id, &req)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteCourseSession handles DELETE /course-session/:id. A
// session of another client is reported exactly like a missing
// one, so ids cannot be probed across tenants.
func (h *Handler) DeleteCourseSession(c *gin.Context) {

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
		"message": "تم حذف محاضرة الدورة بنجاح.",
	})
}