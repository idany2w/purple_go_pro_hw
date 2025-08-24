package verify

import (
	"demo/validation/configs"
	"demo/validation/pkg/response"
	"net/http"
)

type VerifyHandlerDeps struct {
	*configs.Config
	VerifyService *VerifyService
}

type VerifyHandler struct {
	*configs.Config
	VerifyService *VerifyService
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config:        deps.Config,
		VerifyService: deps.VerifyService,
	}

	router.HandleFunc("POST /send", handler.send())
	router.HandleFunc("POST /verify/{hash}", handler.verify())
}

func (h *VerifyHandler) send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := h.VerifyService.GenerateHash(r.FormValue("email"))
		err := h.VerifyService.SendEmail("test@test.com", hash)

		res := VerifySendEmailResponse{
			Status:  "success",
			Error:   "",
			Message: "Email sent",
		}

		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
			return
		}

		response.SendJsonResponse(w, res, http.StatusOK)
	}
}

func (h *VerifyHandler) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isVerified := true

		res := VerifyHashResponse{
			Status:  "success",
			Error:   "",
			Message: "Email verified",
		}

		if isVerified {
			response.SendJsonResponse(w, res, http.StatusOK)
		} else {
			res.Status = "error"
			res.Error = "Email not verified"
			response.SendJsonResponse(w, res, http.StatusOK)
		}
	}
}
