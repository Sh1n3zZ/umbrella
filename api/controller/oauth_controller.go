package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Sh1n3zZ/umbrella/domain"
	"github.com/Sh1n3zZ/umbrella/usecase"
)

// OauthController exposes OAuth 2.0 HTTP handlers (RFC 6749).
type OauthController struct {
	oauthClients *usecase.OauthClientsUsecase
}

// NewOauthController constructs a controller backed by oauth client use cases.
func NewOauthController(oauthClients *usecase.OauthClientsUsecase) *OauthController {
	return &OauthController{oauthClients: oauthClients}
}

// Authorization implements the authorization endpoint for the authorization code grant
// (RFC 6749 Section 4.1.1). Resource owner authentication and consent UI are not implemented;
// valid requests are answered with a successful redirect carrying an authorization code.
//
// GET /v1/oauth/authorize
func (c *OauthController) Authorization(ctx *gin.Context) {
	rctx := ctx.Request.Context()

	req, err := domain.NewAuthorizationRequestFromQuery(
		ctx.Query("client_id"),
		ctx.Query("redirect_uri"),
		ctx.Query("response_type"),
		ctx.Query("state"),
		ctx.Query("scope"),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var clientID pgtype.UUID
	if err := clientID.Scan(req.ClientID.String()); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "client_id must be a UUID"})
		return
	}

	effectiveRedirect, err := c.oauthClients.EffectiveRedirectURI(rctx, clientID, req.RedirectURI)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrOauthUnknownClient),
			errors.Is(err, usecase.ErrOauthRedirectMismatch),
			errors.Is(err, usecase.ErrOauthRedirectRequired):
			ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "failed to process authorization request"})
		}
		return
	}

	if err := domain.ValidateRedirectURI(effectiveRedirect); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if req.ResponseType == "" {
		redir := domain.AuthorizationRedirectError{
			RedirectURI: effectiveRedirect,
			Message:     "response_type is required",
			State:       req.State,
		}
		ctx.Redirect(http.StatusFound, redir.Location())
		return
	}
	if req.ResponseType != domain.ResponseTypeCode {
		redir := domain.AuthorizationRedirectError{
			RedirectURI: effectiveRedirect,
			Message:     "only response_type=code is supported",
			State:       req.State,
		}
		ctx.Redirect(http.StatusFound, redir.Location())
		return
	}

	code, err := c.oauthClients.GenerateURLSafeRandomString(32)
	if err != nil {
		redir := domain.AuthorizationRedirectError{
			RedirectURI: effectiveRedirect,
			Message:     "failed to issue authorization code",
			State:       req.State,
		}
		ctx.Redirect(http.StatusFound, redir.Location())
		return
	}

	success := domain.AuthorizationSuccess{
		RedirectURI: effectiveRedirect,
		Code:        code,
		State:       req.State,
	}
	ctx.Redirect(http.StatusFound, success.Location())
}
