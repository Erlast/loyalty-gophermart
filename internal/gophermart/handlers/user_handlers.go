package handlers

import (
	"encoding/json"
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

func NewUserHandler(
	service *services.UserService,
	logger *zap.SugaredLogger,
) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Register user")
	var user models.User
	if err := render.Bind(r, &user); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	h.logger.Info("User created")

	if err := h.service.Register(r.Context(), &user); err != nil {
		if services.IsDuplicateError(err) {
			h.logger.Error("Username already taken", zap.Error(err))
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		h.logger.Error("Error registering user", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("User registered")

	token, err := services.GenerateJWT(user.ID, h.logger)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Token generated", zap.String("token", token))

	w.Header().Set("Authorization", "Bearer "+token) // Setting the Authorization header with the token
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"Authorization": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := render.Bind(r, &credentials); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		h.logger.Error("Error logging in", zap.Error(err))
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	token, err := services.GenerateJWT(user.ID, h.logger)
	if err != nil {
		h.logger.Error("Error generating JWT", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token) // Setting the Authorization header with the token
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"Authorization": token}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Error writing response", zap.Error(err))
	}
}
