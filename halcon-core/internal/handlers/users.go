package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username   string          `json:"username" validate:"required"`
	Password   string          `json:"password" validate:"required,min=6"`
	Role       models.UserRole `json:"role" validate:"required"`
	Department string          `json:"department"`
	FullName   string          `json:"full_name"`
	Email      string          `json:"email"`
}

type UpdateUserRequest struct {
	Password   string          `json:"password,omitempty"`
	Role       models.UserRole `json:"role"`
	Department string          `json:"department"`
	FullName   string          `json:"full_name"`
	Email      string          `json:"email"`
	IsActive   *bool           `json:"is_active"`
}

// GetUsers returns all users (Admin only)
func GetUsers(c echo.Context) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch users")
	}

	return c.JSON(http.StatusOK, users)
}

// GetUser returns a single user by ID
func GetUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user (Admin only)
func CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Department:   req.Department,
		FullName:     req.FullName,
		Email:        req.Email,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create user")
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user (Admin only)
func UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Update other fields
	user.Role = req.Role
	user.Department = req.Department
	user.FullName = req.FullName
	user.Email = req.Email
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser soft deletes a user (Admin only)
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted successfully"})
}
