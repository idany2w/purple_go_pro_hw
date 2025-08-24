package verify

import (
	"demo/validation/configs"
	"demo/validation/pkg/request"
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
	router.HandleFunc("GET /verify/{hash}", handler.verify())
}

func (h *VerifyHandler) send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := response.Response{}
		body, _ := request.HandleBody[VerifySendEmailRequest](w, r)
		err := h.VerifyService.VerifyEmail(body.Email)

		if err != nil {
			res.Status = response.StatusError
			res.Error = err.Error()
			response.SendJsonResponse(w, res, http.StatusInternalServerError)
			return
		}

		res.Status = response.StatusSuccess
		res.Message = "Email sent"
		response.SendJsonResponse(w, res, http.StatusOK)
	}
}

func (h *VerifyHandler) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := response.Response{}
		hash := r.PathValue("hash")
		isVerified := h.VerifyService.VerifyHash(hash)

		res.Data = VerifyHashRequest{
			IsVerified: isVerified,
		}

		if !isVerified {
			res.Status = response.StatusError
			res.Error = "Is not verified"
			response.SendJsonResponse(w, res, http.StatusNotFound)
			return
		}

		res.Status = response.StatusSuccess
		res.Message = "Email verified"
		response.SendJsonResponse(w, res, http.StatusOK)
	}
}
