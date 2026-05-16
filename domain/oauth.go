package domain

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

const (
	// ResponseTypeCode is the authorization code grant response_type (RFC 6749 Section 4.1.1).
	ResponseTypeCode = "code"
)

// AuthorizationRequest holds query parameters for the authorization endpoint (RFC 6749 Section 4.1.1).
type AuthorizationRequest struct {
	ClientID     uuid.UUID
	RedirectURI  string // optional when exactly one redirect URI is registered
	ResponseType string
	State        string // optional
	Scope        string // optional
}

// NewAuthorizationRequestFromQuery parses authorization endpoint query values into a request.
func NewAuthorizationRequestFromQuery(clientID, redirectURI, responseType, state, scope string) (AuthorizationRequest, error) {
	if clientID == "" {
		return AuthorizationRequest{}, errors.New("client_id is required")
	}
	id, err := uuid.Parse(clientID)
	if err != nil {
		return AuthorizationRequest{}, errors.New("client_id must be a UUID")
	}
	return AuthorizationRequest{
		ClientID:     id,
		RedirectURI:  redirectURI,
		ResponseType: responseType,
		State:        state,
		Scope:        scope,
	}, nil
}

// AuthorizationSuccess is the outcome of a granted authorization (RFC 6749 Section 4.1.2).
type AuthorizationSuccess struct {
	RedirectURI string
	Code        string
	State       string
}

// Location returns the client redirection URI with code and optional state query parameters.
func (s AuthorizationSuccess) Location() string {
	v := url.Values{}
	v.Set("code", s.Code)
	if s.State != "" {
		v.Set("state", s.State)
	}
	loc, _ := MergeRedirectQuery(s.RedirectURI, v)
	return loc
}

// AuthorizationRedirectError is an OAuth error delivered via redirect (RFC 6749 Section 4.1.2.1).
type AuthorizationRedirectError struct {
	RedirectURI string
	Message     string
	State       string
}

// Location returns the client redirection URI with error query parameters.
func (e AuthorizationRedirectError) Location() string {
	v := url.Values{}
	v.Set("error", e.Message)
	if e.State != "" {
		v.Set("state", e.State)
	}
	loc, _ := MergeRedirectQuery(e.RedirectURI, v)
	return loc
}

// ValidateRedirectURI checks RFC 6749 redirection endpoint constraints (Section 3.1.2).
func ValidateRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return errors.New("redirect_uri is not a valid URI")
	}
	if u.Scheme == "" || u.Host == "" {
		return errors.New("redirect_uri must be an absolute URI")
	}
	if u.Fragment != "" {
		return errors.New("redirect_uri must not include a fragment")
	}
	return nil
}

// MergeRedirectQuery merges add into the query component of redirectBase (RFC 3986 Section 3.4).
func MergeRedirectQuery(redirectBase string, add url.Values) (string, error) {
	u, err := url.Parse(redirectBase)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, vals := range add {
		for _, val := range vals {
			q.Set(k, val)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
