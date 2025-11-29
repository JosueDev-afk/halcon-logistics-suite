package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/models"
)

type TrackingRequest struct {
	CustomerNumber string `json:"customer_number" query:"customer_number" validate:"required"`
	InvoiceNumber  string `json:"invoice_number" query:"invoice_number" validate:"required"`
}

type TrackingResponse struct {
	Found              bool                `json:"found"`
	InvoiceNumber      string              `json:"invoice_number,omitempty"`
	CustomerName       string              `json:"customer_name,omitempty"`
	Status             models.OrderStatus  `json:"status,omitempty"`
	DeliveryAddress    string              `json:"delivery_address,omitempty"`
	EvidencePhotoURL   string              `json:"evidence_photo_url,omitempty"`
	CreatedAt          string              `json:"created_at,omitempty"`
	UpdatedAt          string              `json:"updated_at,omitempty"`
}

// TrackOrder allows public tracking of orders without authentication
func TrackOrder(c echo.Context) error {
	var req TrackingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if req.CustomerNumber == "" || req.InvoiceNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "customer_number and invoice_number are required")
	}

	var order models.Order
	err := database.DB.Where(
		"customer_number = ? AND invoice_number = ? AND is_deleted = ?",
		req.CustomerNumber,
		req.InvoiceNumber,
		false,
	).First(&order).Error

	if err != nil {
		return c.JSON(http.StatusOK, TrackingResponse{Found: false})
	}

	return c.JSON(http.StatusOK, TrackingResponse{
		Found:            true,
		InvoiceNumber:    order.InvoiceNumber,
		CustomerName:     order.CustomerName,
		Status:           order.Status,
		DeliveryAddress:  order.DeliveryAddress,
		EvidencePhotoURL: order.EvidencePhotoURL,
		CreatedAt:        order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        order.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}
