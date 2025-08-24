package verify

type VerifyPayload struct {
	Email string `json:"email"`
}

type VerifySendEmailResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type VerifyHashResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
