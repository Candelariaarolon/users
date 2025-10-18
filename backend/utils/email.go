package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"

	log "github.com/sirupsen/logrus"
)

// GenerateVerificationCode generates a random 6-digit code
func GenerateVerificationCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += digit.String()
	}
	return code, nil
}

// SendVerificationEmail sends a verification code to the user's email
func SendVerificationEmail(toEmail, code, userName string) error {
	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("SMTP_FROM")

	// Validate SMTP configuration
	if smtpHost == "" || smtpPort == "" {
		log.Warn("SMTP not configured, skipping email send. Code:", code)
		// In development, just log the code
		fmt.Printf("\n=== VERIFICATION CODE FOR %s ===\n%s\n===============================\n", toEmail, code)
		return nil
	}

	// Create message
	subject := "Email Verification Code"
	body := fmt.Sprintf(`
Hello %s,

Your verification code is: %s

This code will expire in 15 minutes.

If you didn't request this code, please ignore this email.

Best regards,
Users Microservice Team
`, userName, code)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Setup authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, message)
	if err != nil {
		log.Error("Failed to send email:", err)
		// In development, still log the code
		fmt.Printf("\n=== VERIFICATION CODE FOR %s (EMAIL FAILED) ===\n%s\n===============================\n", toEmail, code)
		return err
	}

	log.Info("Verification email sent successfully to:", toEmail)
	return nil
}

// SendWelcomeEmail sends a welcome email after successful verification
func SendWelcomeEmail(toEmail, userName string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("SMTP_FROM")

	if smtpHost == "" || smtpPort == "" {
		log.Warn("SMTP not configured, skipping welcome email")
		return nil
	}

	subject := "Welcome! Your Account is Verified"
	body := fmt.Sprintf(`
Hello %s,

Welcome to our platform! Your email has been successfully verified.

You can now log in and start using our services.

Best regards,
Users Microservice Team
`, userName)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, message)
	if err != nil {
		log.Error("Failed to send welcome email:", err)
		return err
	}

	log.Info("Welcome email sent successfully to:", toEmail)
	return nil
}
