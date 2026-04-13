package dto

import (
	"net/mail"
	"strings"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "is required"
	}
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "is required"
	} else if _, err := mail.ParseAddress(r.Email); err != nil {
		errs["email"] = "is not a valid email"
	}
	if len(r.Password) < 6 {
		errs["password"] = "must be at least 6 characters"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "is required"
	}
	if r.Password == "" {
		errs["password"] = "is required"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserDTO  `json:"user"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
