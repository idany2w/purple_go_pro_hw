package verify

import (
	"crypto/sha256"
	"demo/validation/internal/email"
	"encoding/hex"
)

type VerifyService struct {
	EmailService *email.EmailService
}

func NewVerifyService(emailService *email.EmailService) *VerifyService {
	return &VerifyService{EmailService: emailService}
}

func (s *VerifyService) GenerateHash(email string) string {
	hash := sha256.New()
	hash.Write([]byte(email))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *VerifyService) SendEmail(to string, hash string) error {
	subject := "Verify your email"
	text := "Verify your email"

	return s.EmailService.SendEmail(to, subject, text)
}
