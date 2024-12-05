/**
 *  MockEmailService provides a mock implementation of the EmailServiceInterface for testing purposes.
 *  This mock captures sent emails and allows you to validate email content, recipients, and other details
 *  during testing without actually sending emails.
 *
 *  @struct   MockEmailService
 *  @inherits EmailServiceInterface
 *
 *  @fields
 *  - SentEmails ([]Email): A slice to store the details of emails sent during the test.
 *
 *  @struct   Email
 *  - To (string): The recipient's email address.
 *  - Subject (string): The email subject.
 *  - Body (string): The email body content.
 *
 *  @methods
 *  - SendEmail(toEmail, subject, body) (error): Captures the email details and appends them to the SentEmails slice.
 *
 *  @example
 *  ```
 *  // Initialize the mock email service
 *  mockEmailService := &MockEmailService{}
 *
 *  // Simulate sending an email
 *  err := mockEmailService.SendEmail("test@example.com", "Test Subject", "Test Body")
 *
 *  // Validate that the email was captured
 *  if len(mockEmailService.SentEmails) != 1 {
 *      t.Errorf("Expected 1 email, got %d", len(mockEmailService.SentEmails))
 *  }
 *  email := mockEmailService.SentEmails[0]
 *  fmt.Println(email.To)      // Output: test@example.com
 *  fmt.Println(email.Subject) // Output: Test Subject
 *  fmt.Println(email.Body)    // Output: Test Body
 *  ```
 *
 *  @file      mock_email.go
 *  @project   DailyVerse
 *  @framework Go Testing with Mock Services
 */

package mocks

// MockEmailService is a mock implementation of the EmailServiceInterface.
type MockEmailService struct {
	// SentEmails stores the details of all emails sent during testing.
	SentEmails []Email
}

// Email represents the details of an email sent using the mock service.
type Email struct {
	To      string // Recipient's email address
	Subject string // Email subject
	Body    string // Email body content
}

// SendEmail simulates sending an email by capturing its details.
// Parameters:
// - toEmail (string): Recipient's email address.
// - subject (string): Subject of the email.
// - body (string): Body content of the email.
//
// Returns:
// - error: Always returns nil, as this is a simulation.
func (mes *MockEmailService) SendEmail(toEmail, subject, body string) error {
	// Append the email details to the SentEmails slice.
	mes.SentEmails = append(mes.SentEmails, Email{To: toEmail, Subject: subject, Body: body})
	return nil
}
