package mail

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Sh1n3zZ/umbrella/domain"
)

func validateVerificationMail(msg domain.VerificationMail) error {
	if strings.TrimSpace(msg.To) == "" {
		return errors.New("mail: recipient is required")
	}
	if strings.TrimSpace(msg.Name) == "" {
		return errors.New("mail: name is required")
	}
	if strings.TrimSpace(msg.URL) == "" {
		return errors.New("mail: verify URL is required")
	}
	return nil
}

func (s *smtpSender) SendVerification(ctx context.Context, msg domain.VerificationMail) error {
	if err := validateVerificationMail(msg); err != nil {
		return err
	}

	err := s.sendRendered(ctx, msg.To, "Email Address Verification", func(r *renderer) (string, error) {
		return r.renderVerification(verificationEmail{
			Name: strings.TrimSpace(msg.Name),
			URL:  strings.TrimSpace(msg.URL),
		})
	})
	if err != nil {
		return fmt.Errorf("mail: send verification: %w", err)
	}
	return nil
}
