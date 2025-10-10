package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"massage-booking/backend/models"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig loads email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUser:     getEnvOrDefault("SMTP_USER", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASS", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@massagebooking.com"),
		FromName:     getEnvOrDefault("FROM_NAME", "Massage Booking Team"),
	}
}

// SendConfirmationEmail sends booking confirmation email
func SendConfirmationEmail(booking *models.BookingDetail) error {
	config := GetEmailConfig()

	// If SMTP credentials are not configured, log email instead
	if config.SMTPUser == "" || config.SMTPPassword == "" {
		return logEmailToConsole(booking)
	}

	// Generate email content
	subject := GetEmailSubject(booking)
	body := RenderEmailTemplate(booking)

	// Send email via SMTP
	return sendSMTPEmail(config, booking.Email, subject, body)
}

// logEmailToConsole logs email content to console (fallback when SMTP not configured)
func logEmailToConsole(booking *models.BookingDetail) error {
	subject := GetEmailSubject(booking)
	body := RenderEmailTemplate(booking)

	log.Printf("=== EMAIL NOTIFICATION ===")
	log.Printf("To: %s", booking.Email)
	log.Printf("Subject: %s", subject)
	log.Printf("Body (HTML):\n%s", body)
	log.Printf("=== END EMAIL ===")

	// Also save to file for reference
	filename := fmt.Sprintf("email_%s.html", booking.Reference)
	if err := os.WriteFile(filename, []byte(body), 0644); err != nil {
		log.Printf("Warning: Could not save email to file %s: %v", filename, err)
	} else {
		log.Printf("Email content saved to %s", filename)
	}

	return nil
}

// sendSMTPEmail sends email via SMTP
func sendSMTPEmail(config *EmailConfig, to, subject, htmlBody string) error {
	// SMTP server configuration
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)
	addr := config.SMTPHost + ":" + config.SMTPPort

	// Email headers and body
	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Send email
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		// Fallback to console logging
		return logEmailToConsole(&models.BookingDetail{
			Email:     to,
			Reference: "SMTP_FAILED",
		})
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// SendEmailAsync sends email in background goroutine
func SendEmailAsync(booking *models.BookingDetail) {
	go func() {
		if err := SendConfirmationEmail(booking); err != nil {
			log.Printf("Error sending confirmation email: %v", err)
		}
	}()
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
