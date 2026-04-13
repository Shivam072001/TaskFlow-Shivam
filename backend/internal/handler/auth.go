package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	user, token, err := h.authSvc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteValidationError(w, map[string]string{"email": "already registered"})
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	dto.WriteJSON(w, http.StatusCreated, dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	user, token, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			dto.WriteError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
