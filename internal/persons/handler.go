package persons

import (
	"errors"
	"net/http"
	"strconv"
	"training-app/internal/auth"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	secret  string
}

func NewHandler(service *Service, secret string) *Handler {
	return &Handler{
		service: service,
		secret:  secret,
	}
}

func (h *Handler) ListPersons(c *gin.Context) {

	clientID := c.GetInt("client_id")

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

	clientID := c.GetInt("client_id")

	person, err := h.service.GetByID(clientID, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, person)
}

// ListUsers returns the client's registered users: persons
// who already have credentials. The admin picks a reset
// target from this list, so he cannot select a person
// without a user.
func (h *Handler) ListUsers(c *gin.Context) {

	// client_id was placed into the Gin context
	// by the JWT middleware.
	clientID := c.GetInt("client_id")

	users, err := h.service.GetAllUsers(clientID)

	if err != nil {

		// Anything else is a database error.
		status, message := dbutil.TranslateDBError(err)

		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	c.JSON(http.StatusOK, users)
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

	// client_id and user_id were placed into the Gin context
	// by the JWT middleware. user_id is needed to check which
	// groups the authenticated user may assign.
	clientID := c.GetInt("client_id")
	userID := c.GetInt("user_id")

	person, err := h.service.RegisterUser(clientID, userID, &req)

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

		// The authenticated user does not exist for this client.
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The authenticated user has no group, so he cannot
		// assign anyone.
		case errors.Is(err, ErrCallerHasNoGroup):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "لا يمكن إتمام العملية لأن حسابك غير مرتبط بمجموعة.",
			})

		// The requested group is above the caller's own group.
		case errors.Is(err, ErrGroupNotAllowed):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "لا يمكنك إسناد مستخدم إلى مجموعة أعلى من مجموعتك.",
			})

		// The username became empty after normalization.
		case errors.Is(err, ErrInvalidUsername):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "يرجى إدخال اسم مستخدم صحيح.",
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

// Login verifies the credentials and issues a JWT that the
// client must send as "Authorization: Bearer <token>" on
// every authenticated route.
func (h *Handler) Login(c *gin.Context) {

	var req models.LoginRequest

	// Convert JSON body into our request struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	person, err := h.service.Login(req.Username, req.Password)

	if err != nil {

		// One generic message for wrong username, wrong
		// password or unknown user, so the response does
		// not leak which part failed.
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "اسم المستخدم أو كلمة المرور غير صحيحة.",
			})
			return
		}

		// Anything else is a database error.
		status, message := dbutil.TranslateDBError(err)

		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	// Sign the token with the id of the person and his client,
	// exactly what the JWT middleware later puts back into the
	// Gin context as user_id and client_id.
	token, err := auth.GenerateToken(person.ID, person.ClientID, h.secret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "تعذر إنشاء الجلسة، يرجى المحاولة مرة أخرى.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        person.ID,
			"name":      person.Name,
			"group_id":  person.GroupID,
			"client_id": person.ClientID,
		},
	})
}

// Logout ends the session on the client side. JWTs are
// stateless: the server cannot revoke a token it already
// signed without keeping a denylist, so the client is
// responsible for discarding the token it received from
// /login. The route still requires a valid token (see the
// authenticated route group in cmd/server/main.go), so it
// doubles as a server-side token check.
func (h *Handler) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تسجيل الخروج بنجاح.",
	})
}

// ChangePassword lets the authenticated user change his own
// password after proving he knows the old one.
func (h *Handler) ChangePassword(c *gin.Context) {

	var req models.ChangePasswordRequest

	// Convert JSON body into our request struct.
	// The binding tags enforce the required fields
	// and the minimum password length.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	// client_id and user_id were placed into the Gin context
	// by the JWT middleware. The user can only change his
	// own password, and the id comes from the signed token,
	// never from the request.
	clientID := c.GetInt("client_id")
	userID := c.GetInt("user_id")

	err := h.service.ChangePassword(clientID, userID, &req)

	if err != nil {

		switch {

		// The authenticated user does not exist for this client.
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The person has no credentials yet, so there is
		// no password to change.
		case errors.Is(err, ErrNotRegistered):
			c.JSON(http.StatusConflict, gin.H{
				"error": "الحساب غير مسجل كمستخدم.",
			})

		// The old password did not match the stored hash.
		case errors.Is(err, ErrWrongOldPassword):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "كلمة المرور القديمة غير صحيحة.",
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

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تغيير كلمة المرور بنجاح.",
	})
}

// ResetPassword lets an admin (group 1) set a new password
// for another user without the old one.
func (h *Handler) ResetPassword(c *gin.Context) {

	var req models.ResetPasswordRequest

	// Convert JSON body into our request struct.
	// The binding tags enforce the required fields
	// and the minimum password length.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "بيانات الطلب غير صحيحة.",
		})
		return
	}

	// client_id and user_id were placed into the Gin context
	// by the JWT middleware. user_id is needed to check that
	// the caller is an admin.
	clientID := c.GetInt("client_id")
	userID := c.GetInt("user_id")

	// The target person comes from the URL, not the body.
	targetID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = h.service.ResetPassword(clientID, userID, targetID, &req)

	if err != nil {

		switch {

		// The caller does not exist for this client.
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The target person does not exist for this client.
		case errors.Is(err, ErrPersonNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})

		// The target person has no credentials yet, so
		// RegisterUser is the flow that grants them.
		case errors.Is(err, ErrNotRegistered):
			c.JSON(http.StatusConflict, gin.H{
				"error": "الحساب غير مسجل كمستخدم.",
			})

		// Only admins may reset passwords.
		case errors.Is(err, ErrNotAdmin):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "غير مسموح لك بإعادة تعيين كلمة المرور.",
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

	c.JSON(http.StatusOK, gin.H{
		"message": "تم إعادة تعيين كلمة المرور بنجاح.",
	})
}
