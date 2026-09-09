package courses

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

func (h *Handler) ListCourses(c *gin.Context) {

	clientID := c.GetInt("client_id")

	courses, err := h.service.GetAll(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *Handler) GetCourse(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	clientID := c.GetInt("client_id")

	course, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

// CreateCourse handles POST /courses. It creates a course in
// one of the three valid states: planned-not-executed,
// executed-not-planned or executed-from-planned, and
// enforces that exactly one subject identifier is set.
func (h *Handler) CreateCourse(c *gin.Context) {

	var req models.CreateCourseRequest

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

	course, err := h.service.Create(clientID, &req)

	if err != nil {

		switch {

		// Rule 1: exactly one subject identifier.
		case errors.Is(err, ErrNoSubjectSet),
			errors.Is(err, ErrMultipleSubjectsSet):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "يمكن اختيار خيار واحد فقط لتحديد مادة المقرر.",
			})

		// Rule 2: the lifecycle flags must form a valid state.
		case errors.Is(err, ErrInvalidState):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "حالة المقرر غير صحيحة، يجب أن يكون مخططاً أو منفذاً.",
			})

		case errors.Is(err, ErrUnexpectedPlannedID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "لا يجوز إرسال معرف المقرر المخطط في هذه الحالة.",
			})

		case errors.Is(err, ErrPlannedIDRequired):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "يجب إرسال معرف المقرر المخطط عند تنفيذ مقرر من مقرر مخطط.",
			})

		// The referenced planned course is missing, belongs
		// to another client or was already executed.
		case errors.Is(err, ErrPlannedCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "المقرر المخطط غير موجود أو تم تنفيذه مسبقاً.",
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

	c.JSON(http.StatusCreated, course)
}

// ListPlannedCourses handles GET /get-planned-courses. It
// renders the pending planned courses (is_planned = true and
// is_executed = false) the frontend picks from when executing
// a course from a planned one. The subject query parameters
// are optional filters; when none is sent, all planned
// courses are returned.
func (h *Handler) ListPlannedCourses(c *gin.Context) {

	filters, err := parseCourseListFilters(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clientID := c.GetInt("client_id")

	courses, err := h.service.GetPlanned(clientID, filters)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// parseCourseListFilters reads the optional subject query
// parameters. A parameter that was not sent stays nil so the
// repository skips its WHERE condition; a non-integer value
// is rejected with 400.
func parseCourseListFilters(c *gin.Context) (models.CourseListFilters, error) {

	var filters models.CourseListFilters

	if raw := c.Query("path_grade_subject_id"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return filters, errors.New("invalid path_grade_subject_id")
		}
		filters.PathGradeSubjectID = &value
	}

	if raw := c.Query("path_grade_subject_term_id"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return filters, errors.New("invalid path_grade_subject_term_id")
		}
		filters.PathGradeSubjectTermID = &value
	}

	if raw := c.Query("learning_subject_id"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return filters, errors.New("invalid learning_subject_id")
		}
		filters.LearningSubjectID = &value
	}

	return filters, nil
}
