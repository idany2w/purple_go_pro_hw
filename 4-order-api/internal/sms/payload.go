package sms

type SendSmsRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type SendSmsResponse struct {
	Key string `json:"key"`
}

type VerifySmsRequest struct {
	Phone string `json:"phone" validate:"required"`
	Token string `json:"token" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

type VerifySmsResponse struct {
	Token string `json:"token"`
}
