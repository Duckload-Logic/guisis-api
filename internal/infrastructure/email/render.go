package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

func renderTemplate(email audit.EmailEntry) (string, error) {
	if email.Body != "" {
		return email.Body, nil
	}

	if email.TemplatePath == "" {
		return "", fmt.Errorf("no body or template provided")
	}

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

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, email.TemplateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
