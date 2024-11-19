// tests/mocks/mock_email.go
package mocks

type MockEmailService struct {
	SentEmails []Email
}

type Email struct {
	To      string
	Subject string
	Body    string
}

func (mes *MockEmailService) SendEmail(toEmail, subject, body string) error {
	mes.SentEmails = append(mes.SentEmails, Email{To: toEmail, Subject: subject, Body: body})
	return nil
}
