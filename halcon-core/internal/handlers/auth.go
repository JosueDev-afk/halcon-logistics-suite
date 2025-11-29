package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/models"
	"github.com/nietzshn/halcon-core/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID         uint            `json:"id"`
	Username   string          `json:"username"`
	Role       models.UserRole `json:"role"`
	Department string          `json:"department"`
	FullName   string          `json:"full_name"`
	Email      string          `json:"email"`
}

// Login authenticates a user and returns a JWT token
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Find user by username
	var user models.User
	if err := database.DB.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Role:       user.Role,
			Department: user.Department,
			FullName:   user.FullName,
			Email:      user.Email,
		},
	})
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Role:       user.Role,
		Department: user.Department,
		FullName:   user.FullName,
		Email:      user.Email,
	})
}
