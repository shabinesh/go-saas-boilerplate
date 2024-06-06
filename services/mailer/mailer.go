package mailer

import "fmt"

type MockMailer struct {
}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (m *MockMailer) SendEmail(to, subject, body string) error {
	fmt.Printf("Sending email \nto: %s\nsubject: %s \nbody: \n%s\n", to, subject, body)

	return nil
}
