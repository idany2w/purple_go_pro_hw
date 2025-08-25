package verify

type VerifySendEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyHashRequest struct {
	IsVerified bool `json:"is_verified"`
}
