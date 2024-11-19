// internal/services/email.go
package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type EmailServiceInterface interface {
	SendEmail(toEmail, subject, body string) error
}

type SMTPEmailService struct {
	Auth smtp.Auth
	Host string
	Port int
	From string
}

func NewSMTPEmailService() EmailServiceInterface {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASS"), os.Getenv("SMTP_HOST"))
	return &SMTPEmailService{
		Auth: auth,
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		From: os.Getenv("EMAIL_USER"),
	}
}

func (es *SMTPEmailService) SendEmail(toEmail, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", es.Host, es.Port)
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	return smtp.SendMail(addr, es.Auth, es.From, []string{toEmail}, msg)
}
