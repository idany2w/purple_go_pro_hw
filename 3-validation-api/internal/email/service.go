package email

import (
	"demo/validation/configs"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailService struct {
	*configs.Config
}

func NewEmailService(config *configs.Config) *EmailService {
	return &EmailService{Config: config}
}

func (s *EmailService) SendEmail(to string, subject string, text string) error {
	email := email.NewEmail()
	email.From = s.Config.Email.Email
	email.To = []string{to}
	email.Subject = subject
	email.Text = []byte(text)

	addr := fmt.Sprintf("%s:%s", s.Config.Email.Address, s.Config.Email.Port)
	smtpPlainAuth := smtp.PlainAuth("", s.Config.Email.Email, s.Config.Email.Password, s.Config.Email.Address)

	return email.Send(addr, smtpPlainAuth)
}
