package domain

import "context"

type VerificationMail struct {
	To   string
	Name string
	URL  string
}

type MailSender interface {
	SendVerification(ctx context.Context, msg VerificationMail) error
}
