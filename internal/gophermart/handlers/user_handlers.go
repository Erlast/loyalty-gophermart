package handlers

import (
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *services.UserService
	logger  *zap.SugaredLogger
}

const errorRenderingError = "Error rendering error : %v"

func NewUserHandler(service *services.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Register user")
	var user models.User
	if err := render.Bind(r, &user); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			h.logger.Errorf("Error rendering request: %v", err)
		}
	}
	h.logger.Info("User created")
	if err := h.service.Register(r.Context(), &user); err != nil {
		h.logger.Error("Error registering user", zap.Error(err))
		err := render.Render(w, r, ErrInternalServerError(err))
		if err != nil {
			h.logger.Errorf(errorRenderingError, err)
		}
	}
	h.logger.Info("User registered")
	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		err := render.Render(w, r, ErrInternalServerError(err))
		if err != nil {
			h.logger.Errorf(errorRenderingError, err)
		}
	}
	h.logger.Info("Token generated", zap.String("token", token))

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"Authorization": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := render.Bind(r, &credentials); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			h.logger.Errorf(errorRenderingError, err)
		}
	}

	user, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		h.logger.Error("Error logging in", zap.Error(err))
		err := render.Render(w, r, ErrUnauthorized(err))
		if err != nil {
			h.logger.Errorf(errorRenderingError, err)
		}
	}

	token, err := services.GenerateJWT(user.ID)
	claims, err := services.ParseJWT(token)
	h.logger.Infof("Token generated %v", claims)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		err := render.Render(w, r, ErrInternalServerError(err))
		if err != nil {
			h.logger.Errorf(errorRenderingError, err)
		}
	}

	render.JSON(w, r, map[string]string{"Authorization": token})
}
