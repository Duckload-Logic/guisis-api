package constants

import (
	"fmt"
	"os"
)

func FromEmail() string {
	if from := os.Getenv("SMTP_FROM"); from != "" {
		return from
	}
	sender := "noreply"
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost"
	}
	return fmt.Sprintf("%s@%s", sender, domain)
}
