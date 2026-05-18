package mail

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFS embed.FS

const verificationTemplatePath = "templates/verification.html"

func parseVerificationTemplate() (*template.Template, error) {
	return template.ParseFS(templateFS, verificationTemplatePath)
}
