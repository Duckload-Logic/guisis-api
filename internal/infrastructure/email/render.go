package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

//go:embed templates/*
var templateFS embed.FS

func renderTemplate(email audit.EmailEntry) (string, error) {
	if email.Body != "" {
		return email.Body, nil
	}

	if email.TemplatePath == "" {
		return "", fmt.Errorf("no body or template provided")
	}

	var tmpl *template.Template
	var err error

	if !filepath.IsAbs(email.TemplatePath) {
		// Use embedded templates first
		embedPath := "templates/" + email.TemplatePath
		tmpl, err = template.ParseFS(templateFS, embedPath)
	}

	// Fallback to disk if absolute path or embed failed
	if tmpl == nil || err != nil {
		path := email.TemplatePath
		if !filepath.IsAbs(path) {
			path = filepath.Join(
				"internal",
				"infrastructure",
				"email",
				"templates",
				email.TemplatePath,
			)
		}
		tmpl, err = template.ParseFiles(path)
		if err != nil {
			return "", fmt.Errorf(
				"failed to parse template %s: %w",
				path,
				err,
			)
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, email.TemplateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
