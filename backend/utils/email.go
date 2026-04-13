package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
	fromEmail := os.Getenv("SMTP_FROM")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}

	subject := "Email Verification Code"
	body := fmt.Sprintf(`Hello %s,

Your verification code is: %s

This code will expire in 15 minutes.

If you didn't request this code, please ignore this email.

Best regards,
UniChat Team`, userName, code)

	err := sendResendEmail(fromEmail, toEmail, subject, body)
	if err != nil {
		log.Error("Failed to send verification email:", err)
		fmt.Printf("\n=== VERIFICATION CODE FOR %s (EMAIL FAILED) ===\n%s\n===============================\n", toEmail, code)
		return err
	}

	log.Info("Verification email sent successfully to:", toEmail)
	return nil
}

// SendWelcomeEmail sends a welcome email after successful verification
func SendWelcomeEmail(toEmail, userName string) error {
	fromEmail := os.Getenv("SMTP_FROM")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}

	subject := "Welcome! Your Account is Verified"
	body := fmt.Sprintf(`Hola %s,

¡Bienvenido a nuestra plataforma! Tu correo electrónico ha sido verificado con éxito.

Ahora puedes iniciar sesión y comenzar a utilizar UniChat.`, userName)

	err := sendResendEmail(fromEmail, toEmail, subject, body)
	if err != nil {
		log.Error("Failed to send welcome email:", err)
		return err
	}

	log.Info("Welcome email sent successfully to:", toEmail)
	return nil
}

// sendResendEmail sends an email via Resend's HTTP API.
// Requires SMTP_PASS to be set to the Resend API key.
func sendResendEmail(from, to, subject, text string) error {
	apiKey := os.Getenv("SMTP_PASS")
	if apiKey == "" {
		log.Warn("SMTP_PASS (Resend API key) not configured, skipping email")
		return nil
	}

	payload := map[string]interface{}{
		"from":    from,
		"to":      []string{to},
		"subject": subject,
		"text":    text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var resBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&resBody)
		return fmt.Errorf("resend API error %d: %v", resp.StatusCode, resBody)
	}

	return nil
}
