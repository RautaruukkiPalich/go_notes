package server

type (
	errorResponse struct {
		Code int `json:"code"`
		Error string `json:"error"`
	}
)