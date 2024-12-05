/**
 *  Email Service provides functionality to send emails using the SMTP protocol.
 *  It leverages environment variables for configuration, ensuring secure and flexible setup.
 *
 *  @interface EmailServiceInterface
 *  @struct   SMTPEmailService
 *  @methods
 *  - NewSMTPEmailService()         - Initializes a new SMTPEmailService instance with environment configurations.
 *  - SendEmail(toEmail, subject, body) - Sends an email to the specified recipient.
 *
 *  @dependencies
 *  - net/smtp: Provides the SMTP client for sending emails.
 *  - os.Getenv: Fetches configuration values from environment variables.
 *  - strconv.Atoi: Converts port string to an integer.
 *
 *  @file      email.go
 *  @project   DailyVerse
 *  @purpose   Utility service for email communication in the application.
 *  @framework Go Standard Library with SMTP Integration
 *  @environment_variables
 *  - SMTP_HOST: The hostname of the SMTP server (e.g., smtp.gmail.com).
 *  - SMTP_PORT: The port number of the SMTP server (e.g., 587).
 *  - EMAIL_USER: The email address used to send emails.
 *  - EMAIL_PASS: The password or app-specific password for the sending email account.
 *
 *  @example
 *  ```
 *  emailService := NewSMTPEmailService()
 *  err := emailService.SendEmail("recipient@example.com", "Welcome to DailyVerse", "Thank you for joining!")
 *  if err != nil {
 *      log.Fatalf("Failed to send email: %v", err)
 *  }
 *  ```
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

// EmailServiceInterface defines the contract for email services.
type EmailServiceInterface interface {
	// SendEmail sends an email with the specified subject and body to the recipient.
	SendEmail(toEmail, subject, body string) error
}

// SMTPEmailService implements EmailServiceInterface using the SMTP protocol.
type SMTPEmailService struct {
	Auth smtp.Auth // Authentication credentials for the SMTP server.
	Host string    // SMTP server hostname.
	Port int       // SMTP server port number.
	From string    // Sender's email address.
}

// NewSMTPEmailService initializes an SMTPEmailService using environment variables for configuration.
// Required environment variables:
// - SMTP_HOST: SMTP server hostname.
// - SMTP_PORT: SMTP server port.
// - EMAIL_USER: Email address used for sending.
// - EMAIL_PASS: Password for the email address.
func NewSMTPEmailService() EmailServiceInterface {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) // Convert port to integer.
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASS"), os.Getenv("SMTP_HOST"))
	return &SMTPEmailService{
		Auth: auth,
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		From: os.Getenv("EMAIL_USER"),
	}
}

// SendEmail sends an email using the SMTP server.
// Parameters:
// - toEmail (string): Recipient's email address.
// - subject (string): Email subject.
// - body (string): Email body.
// Returns:
// - error: Returns an error if the email cannot be sent.
func (es *SMTPEmailService) SendEmail(toEmail, subject, body string) error {
	// Construct the SMTP server address.
	addr := fmt.Sprintf("%s:%d", es.Host, es.Port)

	// Create the email message.
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email using the configured SMTP server.
	return smtp.SendMail(addr, es.Auth, es.From, []string{toEmail}, msg)
}
