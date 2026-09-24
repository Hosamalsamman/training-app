package pathGradeCandidate

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

func (h *Handler) ListPathGradeCandidates(c *gin.Context) {

	clientID := c.GetInt("client_id")

	candidates, err := h.service.GetAll(clientID)

	if err != nil {

		log.Printf("ListPathGradeCandidates failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, candidates)
}

func (h *Handler) GetPathGradeCandidate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	clientID := c.GetInt("client_id")

	candidate, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, candidate)
}

// translateError maps the sentinel errors of the service to
// the HTTP status code and Arabic message each rule deserves.
// Anything else falls through to the DB error translator,
// exactly like the course session handlers do.
func translateError(err error) (int, gin.H) {

	switch {

	// The candidate id is missing, or belongs to another
	// client: both are reported the same way so ids cannot
	// be probed across tenants.
	case errors.Is(err, ErrPathGradeCandidateNotFound):
		return http.StatusNotFound, gin.H{
			"error": "المرشح غير موجود.",
		}

	// The person referenced by person_id is missing, or
	// belongs to another client.
	case errors.Is(err, ErrPersonNotFound):
		return http.StatusNotFound, gin.H{
			"error": "الموظف غير موجود.",
		}

	// The path grade referenced by path_grade_id is
	// missing, or belongs to another client.
	case errors.Is(err, ErrPathGradeNotFound):
		return http.StatusNotFound, gin.H{
			"error": "درجة المسار غير موجودة.",
		}

	// The person has no current grade, so the grade rules
	// cannot be evaluated.
	case errors.Is(err, ErrPersonHasNoCurrentGrade):
		return http.StatusBadRequest, gin.H{
			"error": "الموظف ليس لديه درجة حالية.",
		}

	// The person has no date for the current grade, so the
	// interval rule cannot be evaluated.
	case errors.Is(err, ErrPersonHasNoGradeDate):
		return http.StatusBadRequest, gin.H{
			"error": "الموظف ليس لديه تاريخ للدرجة الحالية.",
		}

	// The current or target grade row is missing its
	// numeric code.
	case errors.Is(err, ErrGradeCodeMissing):
		return http.StatusBadRequest, gin.H{
			"error": "الدرجة لا تحتوي على رقم.",
		}

	// The current grade row is missing its interval.
	case errors.Is(err, ErrGradeIntervalMissing):
		return http.StatusBadRequest, gin.H{
			"error": "الدرجة الحالية لا تحتوي على مدة محددة.",
		}

	// The next-grade rule: the target grade must be exactly
	// one step above the current grade.
	case errors.Is(err, ErrGradeNotNext):
		return http.StatusBadRequest, gin.H{
			"error": "الدرجة المختارة يجب أن تكون الدرجة التالية مباشرة للدرجة الحالية.",
		}

	// The interval rule: the person has not served the
	// interval of the current grade yet.
	case errors.Is(err, ErrGradeIntervalNotServed):
		return http.StatusBadRequest, gin.H{
			"error": "لم يمر على الموظف الوقت المحدد للترقية من درجته الحالية.",
		}

	// Anything else is a database error.
	default:
		status, message := dbutil.TranslateDBError(err)

		return status, gin.H{
			"error": message,
		}
	}
}

// Create handles POST /new-path-grade-candidate.
func (h *Handler) Create(c *gin.Context) {

	var req models.PathGradeCandidate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات غير صحيحة",
		})
		return
	}

	// client_id was placed into the Gin context
	// by the JWT middleware.
	clientID := c.GetInt("client_id")

	_, err := h.service.Create(clientID, &req) // if you need to return the created object use res instead of _ and replace gin.H with res

	if err != nil {

		status, msg := dbutil.TranslateDBError(err)

		c.JSON(status, gin.H{
			"error": msg,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": "تمت الإضافة بنجاح",
	})
}

// Update handles PUT /path-grade-candidate/:id. It is a full
// update of the candidate foreign keys.
func (h *Handler) Update(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req models.PathGradeCandidate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات غير صحيحة",
		})
		return
	}

	clientID := c.GetInt("client_id")

	candidate, err := h.service.Update(clientID, id, &req)

	if err != nil {
		status, body := translateError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, candidate)
}

// Delete handles DELETE /path-grade-candidate/:id. A candidate
// of another client is reported exactly like a missing one, so
// ids cannot be probed across tenants.
func (h *Handler) Delete(c *gin.Context) {

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
		"message": "تم حذف المرشح بنجاح.",
	})
}
