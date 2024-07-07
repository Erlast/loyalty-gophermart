package handlers

import (
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/services"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *services.UserService
	logger  *zap.SugaredLogger
}

func NewUserHandler(service *services.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := render.Bind(r, &user); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.service.Register(r.Context(), &user); err != nil {
		h.logger.Error("Error registering user", zap.Error(err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"Authorization": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := render.Bind(r, &credentials); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		h.logger.Error("Error logging in", zap.Error(err))
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, map[string]string{"Authorization": token})
}
