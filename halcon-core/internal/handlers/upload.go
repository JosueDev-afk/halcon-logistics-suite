package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/config"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/models"
)

type UploadResponse struct {
	URL string `json:"url"`
}

// UploadEvidence handles photo evidence upload for orders
func UploadEvidence(c echo.Context) error {
	orderID := c.Param("id")
	userRole := c.Get("role").(models.UserRole)

	// Only Route role can upload evidence
	if userRole != models.RoleRoute {
		return echo.NewHTTPError(http.StatusForbidden, "only route personnel can upload evidence")
	}

	// Get the order
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	// Get file from request
	file, err := c.FormFile("photo")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "photo file is required")
	}

	// Validate file size
	if file.Size > config.AppConfig.MaxUploadSize {
		return echo.NewHTTPError(http.StatusBadRequest, "file size exceeds maximum allowed")
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return echo.NewHTTPError(http.StatusBadRequest, "only JPG and PNG files are allowed")
	}

	// Create uploads directory if it doesn't exist
	uploadDir := config.AppConfig.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create upload directory")
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("order_%s_%d%s", orderID, timestamp, ext)
	filepath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create file")
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file")
	}

	// Update order with photo URL
	photoURL := fmt.Sprintf("/uploads/%s", filename)
	order.EvidencePhotoURL = photoURL

	// If status is being changed to Delivered, update it
	newStatus := c.FormValue("status")
	if newStatus == string(models.StatusDelivered) {
		order.Status = models.StatusDelivered
		userID := c.Get("user_id").(uint)
		order.LastModifiedBy = userID
	}

	if err := database.DB.Save(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update order")
	}

	return c.JSON(http.StatusOK, UploadResponse{
		URL: photoURL,
	})
}
