package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Sh1n3zZ/umbrella/domain"
)

// OauthUsersController exposes OAuth user account HTTP handlers.
type OauthUsersController struct {
	oauthUsers domain.OauthUsersUsecase
}

// NewOauthUsersController constructs a controller backed by oauth user use cases.
func NewOauthUsersController(oauthUsers domain.OauthUsersUsecase) *OauthUsersController {
	return &OauthUsersController{oauthUsers: oauthUsers}
}

type registerUserRequest struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type verifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// Register creates a new oauth_users row after validating email and hashing the password.
//
// POST /v1/oauth/users/register
func (c *OauthUsersController) Register(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: registerBindErrorMessage(err)})
		return
	}

	resp, err := c.oauthUsers.RegisterOauthUser(
		ctx.Request.Context(),
		req.Nickname,
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, domain.ErrOauthUserEmailTaken) {
			ctx.JSON(http.StatusConflict, domain.ErrorResponse{Message: err.Error()})
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "failed to register user"})
			return
		}
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// Verify finalises the user registration flow by consuming a single-use
// verification token and marking the matching account as email-verified.
//
// POST /v1/users/verify
func (c *OauthUsersController) Verify(ctx *gin.Context) {
	var req verifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "token is required"})
		return
	}

	resp, err := c.oauthUsers.VerifyEmail(ctx.Request.Context(), req.Token)
	if err != nil {
		if errors.Is(err, domain.ErrEmailVerificationTokenInvalid) {
			ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "failed to verify user"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func registerBindErrorMessage(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "invalid request body"
	}
	for _, fe := range ve {
		switch fe.Field() {
		case "Email":
			if fe.Tag() == "email" {
				return "email must be a valid email address"
			}
			return "email is required"
		case "Nickname":
			switch fe.Tag() {
			case "max":
				return "nickname must be at most 64 characters"
			case "min":
				return "nickname is required"
			default:
				return "nickname is required"
			}
		case "Password":
			switch fe.Tag() {
			case "min":
				return "password must be at least 8 characters"
			case "max":
				return "password must be at most 72 characters"
			default:
				return "password is required"
			}
		}
	}
	return "invalid request body"
}
