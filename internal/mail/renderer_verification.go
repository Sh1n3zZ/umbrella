package mail

import (
	"bytes"
	"fmt"
)

type verificationEmail struct {
	Name string
	URL  string
}

func (r *renderer) renderVerification(data verificationEmail) (html string, err error) {
	if r == nil || r.verification == nil {
		return "", fmt.Errorf("mail: verification template is not loaded")
	}

	var buf bytes.Buffer
	if err := r.verification.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("mail: execute verification template: %w", err)
	}
	return buf.String(), nil
}
