package mail

import (
	"fmt"
	"html/template"
)

type renderer struct {
	verification *template.Template
}

func newRenderer() (*renderer, error) {
	verification, err := parseVerificationTemplate()
	if err != nil {
		return nil, fmt.Errorf("mail: load templates: %w", err)
	}
	return &renderer{verification: verification}, nil
}
