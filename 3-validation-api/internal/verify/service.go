package verify

import (
	"crypto/sha256"
	"demo/validation/configs"
	"demo/validation/internal/email"
	"encoding/hex"
	"fmt"
)

const (
	hashSalt = "123123123123"
)

type HashStorage interface {
	SaveHash(hash string) error
	RemoveHash(hash string) error
	GetHash(hash string) (string, error)
}

type VerifyService struct {
	Config       *configs.Config
	EmailService *email.EmailService
	HashStorage  HashStorage
}

func NewVerifyService(config *configs.Config, emailService *email.EmailService, hashStorage HashStorage) *VerifyService {
	return &VerifyService{
		Config:       config,
		EmailService: emailService,
		HashStorage:  hashStorage,
	}
}

func (s *VerifyService) VerifyEmail(email string) error {
	hash, err := s.GenerateHash(email)

	if err != nil {
		return err
	}

	err = s.SaveHash(hash)

	if err != nil {
		return err
	}

	return s.SendEmail(email, hash)
}

func (s *VerifyService) GenerateHash(email string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(fmt.Sprintf("%s%s", email, hashSalt)))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *VerifyService) SaveHash(hash string) error {
	return s.HashStorage.SaveHash(hash)
}

func (s *VerifyService) RemoveHash(hash string) error {
	return s.HashStorage.RemoveHash(hash)
}

func (s *VerifyService) GetHash(hash string) (string, error) {
	return s.HashStorage.GetHash(hash)
}

func (s *VerifyService) SendEmail(to string, hash string) error {
	subject := "Verify your email"
	text := fmt.Sprintf("http://%s:%s/verify/%s", s.Config.Server.Host, s.Config.Server.Port, hash)

	return s.EmailService.SendEmail(to, subject, text)
}

func (s *VerifyService) VerifyHash(hash string) bool {
	content, err := s.GetHash(hash)

	if err != nil {
		return false
	}

	isVerified := content == hash

	if isVerified {
		s.RemoveHash(hash)
	}

	return isVerified
}
