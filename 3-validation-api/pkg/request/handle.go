package request

import (
	"demo/validation/pkg/response"
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	res := response.Response{}

	body, err := Decode[T](r.Body)

	if err != nil {
		res.Status = "error"
		res.Error = err.Error()
		response.SendJsonResponse(w, res, http.StatusBadRequest)
		return nil, err
	}

	err = IsValid(body)

	if err != nil {
		res.Status = "error"
		res.Error = err.Error()
		response.SendJsonResponse(w, res, http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
