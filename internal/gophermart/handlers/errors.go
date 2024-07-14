package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	InvalidOrderFormatMsg = "Invalid order number format"
)

type ErrResponse struct {
	Err            error  `json:"-"`               // низкоуровневая ошибка
	StatusText     string `json:"status"`          // статус ошибки
	ErrorText      string `json:"error,omitempty"` // сообщение об ошибке
	AppCode        int64  `json:"code,omitempty"`  // приложение-специфичный код ошибки
	HTTPStatusCode int    `json:"-"`               // HTTP статус-код
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
