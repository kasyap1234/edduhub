package handler

import (
	"eduhub/server/internal/services/auth"
    "net/http"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
    authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService)*AuthHandler{
    return &AuthHandler{
        authService: authService,
    }
}

type RegisterRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    CollegeName string `json:"collegeName"`
}

func(h *AuthHandler)RegisterUser(e echo.Context)error {
    var req RegisterRequest
    if err :=e.Bind(&req); err !=nil {
        return e.JSON(http.StatusBadRequest,map[string]string{"error": "invalid request body"})

    }
    params := auth.RegisterUserParams
}