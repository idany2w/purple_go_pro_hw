package sms

import (
	"errors"
	"math/rand"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type SmsService struct {
	db *gorm.DB
}

func NewSmsService(db *gorm.DB) *SmsService {
	return &SmsService{db: db}
}

func (s *SmsService) SendSms(phone string) (string, error) {
	phone = s.CorrectPhone(phone)

	if !s.IsValidPhone(phone) {
		return "", errors.New("invalid phone number")
	}

	var existingSms Sms
	err := s.db.Where("phone = ?", phone).First(&existingSms).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.New("failed to send SMS")
	}

	if (err == gorm.ErrRecordNotFound || existingSms.ID != 0) && time.Since(existingSms.CreatedAt) < smsTTL {
		return "", errors.New("SMS already sent. Try again later")
	}

	if err == gorm.ErrRecordNotFound || time.Since(existingSms.CreatedAt) > smsTTL {
		s.db.Delete(&existingSms)
	}

	token := RandStringRunes(TokenLength)
	code := RandStringRunes(CodeLength)
	s.db.Create(&Sms{Phone: phone, Token: token, Code: code})

	// send sms to phone

	return token, nil
}

func (s *SmsService) VerifySms(phone string, token string, code string) error {
	var sms Sms
	phone = s.CorrectPhone(phone)

	err := s.db.Where("phone = ? AND token = ? AND code = ?", phone, token, code).First(&sms).Error
	if err != nil {
		return errors.New("invalid code")
	}

	s.db.Delete(&sms)

	return nil
}

func (s *SmsService) CorrectPhone(phone string) string {
	// Удаляем все символы кроме цифр и +
	re := regexp.MustCompile(`[^\d+]`)
	cleaned := re.ReplaceAllString(phone, "")

	// Если номер начинается с 8, заменяем на +7
	if len(cleaned) > 0 && cleaned[0] == '8' {
		cleaned = "+7" + cleaned[1:]
	}

	// Если номер начинается с 7 без +, добавляем +
	if len(cleaned) > 0 && cleaned[0] == '7' && len(cleaned) == 11 {
		cleaned = "+" + cleaned
	}

	// Если номер не начинается с +, добавляем +7
	if len(cleaned) > 0 && cleaned[0] != '+' {
		if len(cleaned) == 10 {
			cleaned = "+7" + cleaned
		} else if len(cleaned) == 11 && cleaned[0] == '7' {
			cleaned = "+" + cleaned
		}
	}

	return cleaned
}

func (s *SmsService) IsValidPhone(phone string) bool {
	// Проверяем российский номер: +7XXXXXXXXXX (12 символов)
	re := regexp.MustCompile(`^\+7\d{10}$`)
	return re.MatchString(phone)
}

const smsTTL = 10 * time.Minute

func (s *SmsService) DeleteOldSms() {
	s.db.Where("created_at < ?", time.Now().Add(-1*smsTTL)).Delete(&Sms{})
}

const (
	CodeLength  = 6
	TokenLength = 32
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
