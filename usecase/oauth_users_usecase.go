package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Sh1n3zZ/umbrella/domain"
	"github.com/Sh1n3zZ/umbrella/internal/tokenutil"
)

// OauthUsersUsecase coordinates oauth_users-related business rules, including
// the email-verification sub-flow of registration.
type OauthUsersUsecase struct {
	oauthUsersRepo    domain.OauthUsersRepository
	mailSender        domain.MailSender
	verificationStore domain.EmailVerificationStore
	verifyURLTemplate string
}

var _ domain.OauthUsersUsecase = (*OauthUsersUsecase)(nil)

// NewOauthUsersUsecase builds a usecase backed by the given repository, mail
// sender, verification store, and verify-URL template (fmt-style, expects a
// single %s placeholder for the token).
func NewOauthUsersUsecase(
	oauthUsersRepo domain.OauthUsersRepository,
	mailSender domain.MailSender,
	verificationStore domain.EmailVerificationStore,
	verifyURLTemplate string,
) *OauthUsersUsecase {
	return &OauthUsersUsecase{
		oauthUsersRepo:    oauthUsersRepo,
		mailSender:        mailSender,
		verificationStore: verificationStore,
		verifyURLTemplate: verifyURLTemplate,
	}
}

// RegisterOauthUser hashes the password, persists a new user, then issues a
// single-use verification token and dispatches the verification email. The
// token / mail steps are best-effort: failures are logged but do not roll
// back the created account.
func (u *OauthUsersUsecase) RegisterOauthUser(
	ctx context.Context,
	nickname,
	email,
	password string,
) (domain.UserResponse, error) {
	nickname = strings.TrimSpace(nickname)
	email = strings.TrimSpace(strings.ToLower(email))

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.UserResponse{}, err
	}

	user, err := u.oauthUsersRepo.CreateOauthUser(ctx, nickname, email, string(hashed))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserResponse{}, domain.ErrOauthUserEmailTaken
		}
		return domain.UserResponse{}, err
	}

	u.dispatchVerificationEmail(ctx, user)
	return toUserResponse(user), nil
}

// VerifyEmail consumes a verification token and marks the matching user as
// email-verified. Invalid, expired, or already-used tokens map to
// domain.ErrEmailVerificationTokenInvalid.
func (u *OauthUsersUsecase) VerifyEmail(
	ctx context.Context,
	token string,
) (domain.UserResponse, error) {
	if u.verificationStore == nil {
		return domain.UserResponse{}, errors.New("verification store is not configured")
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return domain.UserResponse{}, domain.ErrEmailVerificationTokenInvalid
	}

	userID, err := u.verificationStore.Consume(ctx, token)
	if err != nil {
		return domain.UserResponse{}, err
	}

	user, err := u.oauthUsersRepo.MarkOauthUserEmailVerified(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserResponse{}, domain.ErrEmailVerificationTokenInvalid
		}
		return domain.UserResponse{}, err
	}
	return toUserResponse(user), nil
}

// dispatchVerificationEmail issues a verification token, stores it in Redis,
// and sends the verification email. Failures are logged and swallowed so
// registration still returns success.
func (u *OauthUsersUsecase) dispatchVerificationEmail(ctx context.Context, user domain.User) {
	if u.verificationStore == nil || u.mailSender == nil {
		return
	}

	token, err := tokenutil.GenerateURLSafeRandomString(domain.EmailVerificationTokenBytes)
	if err != nil {
		log.Printf("oauth_users: generate verification token for %s: %v", user.UserID, err)
		return
	}

	if err := u.verificationStore.Save(ctx, token, user.UserID, domain.EmailVerificationTTL); err != nil {
		log.Printf("oauth_users: persist verification token for %s: %v", user.UserID, err)
		return
	}

	verifyURL := formatTokenURL(u.verifyURLTemplate, token)
	if err := u.mailSender.SendVerification(ctx, domain.VerificationMail{
		To:   user.Email,
		Name: user.Nickname,
		URL:  verifyURL,
	}); err != nil {
		log.Printf("oauth_users: send verification email to %s: %v", user.Email, err)
	}
}

// formatTokenURL renders a verification URL by substituting a token into
// template. If template contains "%s" it is treated as a fmt-style pattern;
// otherwise it is treated as a base URL and the token is appended as a
// "token=<value>" query parameter. An empty template falls back to the raw
// token.
func formatTokenURL(template, token string) string {
	if template == "" {
		return token
	}
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, token)
	}
	separator := "?"
	if strings.Contains(template, "?") {
		separator = "&"
	}
	return template + separator + "token=" + token
}

// userResponseFrom maps a User to UserResponse, stripping sensitive fields.
func toUserResponse(u domain.User) domain.UserResponse {
	return domain.UserResponse{
		UserID:        u.UserID.String(),
		Nickname:      u.Nickname,
		Email:         u.Email,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt.UTC().Format(time.RFC3339),
	}
}
