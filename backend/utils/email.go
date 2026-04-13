package utils

import (
	"crypto/rand"
	"crypto/tls"
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

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := sendMail(addr, smtpHost, smtpUser, smtpPass, fromEmail, toEmail, message)
	if err != nil {
		log.Error("Failed to send email:", err)
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
Hola %s,

¡Bienvenido a nuestra plataforma! Tu correo electrónico ha sido verificado con éxito.

Ahora puedes iniciar sesión y comenzar a utilizar UniChat.
`, userName)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := sendMail(addr, smtpHost, smtpUser, smtpPass, fromEmail, toEmail, message)
	if err != nil {
		log.Error("Failed to send welcome email:", err)
		return err
	}

	log.Info("Welcome email sent successfully to:", toEmail)
	return nil
}

// sendMail supports both port 465 (implicit TLS) and port 587 (STARTTLS).
// Go's smtp.SendMail only handles STARTTLS, so port 465 requires a manual TLS dial.
func sendMail(addr, host, user, pass, from, to string, message []byte) error {
	auth := smtp.PlainAuth("", user, pass, host)

	// Port 465 uses implicit TLS — dial TLS first, then speak SMTP.
	if len(addr) >= 3 && addr[len(addr)-3:] == "465" {
		tlsConfig := &tls.Config{ServerName: host}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("tls dial: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("smtp client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err = client.Mail(from); err != nil {
			return fmt.Errorf("smtp mail from: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data: %w", err)
		}
		if _, err = w.Write(message); err != nil {
			return fmt.Errorf("smtp write: %w", err)
		}
		return w.Close()
	}

	// Port 587 (or any other): standard STARTTLS via smtp.SendMail.
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}
