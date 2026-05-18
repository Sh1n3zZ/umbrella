package mail

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Sh1n3zZ/umbrella/domain"
	gomail "github.com/wneessen/go-mail"
)

type smtpSender struct {
	client   *gomail.Client
	from     string
	renderer *renderer
}

func NewSMTPSender(client *gomail.Client, from string) (domain.MailSender, error) {
	from = strings.TrimSpace(from)
	if client == nil {
		return nil, errors.New("mail: smtp client is required")
	}
	if from == "" {
		return nil, errors.New("mail: from address is required")
	}

	r, err := newRenderer()
	if err != nil {
		return nil, err
	}
	return &smtpSender{
		client:   client,
		from:     from,
		renderer: r,
	}, nil
}

var _ domain.MailSender = (*smtpSender)(nil)

func (s *smtpSender) validateReady() error {
	if s == nil {
		return errors.New("mail: sender is nil")
	}
	if s.client == nil {
		return errors.New("mail: smtp client is nil")
	}
	if s.renderer == nil {
		return errors.New("mail: renderer is not initialized")
	}
	if strings.TrimSpace(s.from) == "" {
		return errors.New("mail: from address is required")
	}
	return nil
}

func (s *smtpSender) sendRendered(
	ctx context.Context,
	to, subject string,
	render func(*renderer) (string, error),
) error {
	if err := s.validateReady(); err != nil {
		return err
	}

	htmlBody, err := render(s.renderer)
	if err != nil {
		return err
	}
	return s.sendHTML(ctx, to, subject, htmlBody)
}

func (s *smtpSender) sendHTML(ctx context.Context, to, subject, htmlBody string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("mail: recipient is required")
	}
	if strings.TrimSpace(subject) == "" {
		return errors.New("mail: subject is required")
	}
	if strings.TrimSpace(htmlBody) == "" {
		return errors.New("mail: body is required")
	}

	m := gomail.NewMsg()
	if err := m.From(s.from); err != nil {
		return fmt.Errorf("mail: set from: %w", err)
	}
	if err := m.AddTo(to); err != nil {
		return fmt.Errorf("mail: set to: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(gomail.TypeTextHTML, htmlBody)

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("mail: deliver message: %w", err)
	}
	return nil
}
