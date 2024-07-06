package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // низкоуровневая ошибка
	HTTPStatusCode int   `json:"-"` // HTTP статус-код

	StatusText string `json:"status"`          // статус ошибки
	AppCode    int64  `json:"code,omitempty"`  // приложение-специфичный код ошибки
	ErrorText  string `json:"error,omitempty"` // сообщение об ошибке
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
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
