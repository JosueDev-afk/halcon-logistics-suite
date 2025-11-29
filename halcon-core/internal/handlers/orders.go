package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nietzshn/halcon-core/internal/database"
	"github.com/nietzshn/halcon-core/internal/models"
)

type CreateOrderRequest struct {
	InvoiceNumber   string `json:"invoice_number" validate:"required"`
	CustomerName    string `json:"customer_name" validate:"required"`
	CustomerNumber  string `json:"customer_number" validate:"required"`
	DeliveryAddress string `json:"delivery_address"`
	Notes           string `json:"notes"`
}

type UpdateOrderRequest struct {
	Status          models.OrderStatus `json:"status"`
	DeliveryAddress string             `json:"delivery_address"`
	Notes           string             `json:"notes"`
}

type OrderFilter struct {
	InvoiceNumber  string
	CustomerName   string
	CustomerNumber string
	Status         string
	IncludeDeleted bool
}

// GetOrders returns all orders with optional filters
func GetOrders(c echo.Context) error {
	userRole := c.Get("role").(models.UserRole)

	filter := OrderFilter{
		InvoiceNumber:  c.QueryParam("invoice_number"),
		CustomerName:   c.QueryParam("customer_name"),
		CustomerNumber: c.QueryParam("customer_number"),
		Status:         c.QueryParam("status"),
		IncludeDeleted: c.QueryParam("include_deleted") == "true",
	}

	query := database.DB.Preload("CreatedByUser").Preload("LastModifiedUser")

	// Apply filters
	if filter.InvoiceNumber != "" {
		query = query.Where("invoice_number ILIKE ?", "%"+filter.InvoiceNumber+"%")
	}
	if filter.CustomerName != "" {
		query = query.Where("customer_name ILIKE ?", "%"+filter.CustomerName+"%")
	}
	if filter.CustomerNumber != "" {
		query = query.Where("customer_number ILIKE ?", "%"+filter.CustomerNumber+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Handle soft deletes
	if filter.IncludeDeleted {
		query = query.Unscoped().Where("is_deleted = ?", true)
	} else {
		query = query.Where("is_deleted = ?", false)
	}

	// Role-based filtering
	if userRole == models.RolePurchasing {
		// Purchasing can only see orders in process
		query = query.Where("status = ?", models.StatusInProcess)
	}

	var orders []models.Order
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch orders")
	}

	return c.JSON(http.StatusOK, orders)
}

// GetOrder returns a single order by ID
func GetOrder(c echo.Context) error {
	id := c.Param("id")

	var order models.Order
	query := database.DB.Preload("CreatedByUser").Preload("LastModifiedUser")

	if err := query.First(&order, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	return c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order (Sales only)
func CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	order := models.Order{
		InvoiceNumber:   req.InvoiceNumber,
		CustomerName:    req.CustomerName,
		CustomerNumber:  req.CustomerNumber,
		DeliveryAddress: req.DeliveryAddress,
		Notes:           req.Notes,
		Status:          models.StatusOrdered,
		CreatedBy:       userID,
		LastModifiedBy:  userID,
		IsDeleted:       false,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create order")
	}

	// Reload with associations
	database.DB.Preload("CreatedByUser").Preload("LastModifiedUser").First(&order, order.ID)

	return c.JSON(http.StatusCreated, order)
}

// UpdateOrder updates an existing order
func UpdateOrder(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)
	userRole := c.Get("role").(models.UserRole)

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	var req UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Role-based status transition validation
	if req.Status != "" && req.Status != order.Status {
		if err := validateStatusTransition(order.Status, req.Status, userRole); err != nil {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		order.Status = req.Status
	}

	// Update other fields
	if req.DeliveryAddress != "" {
		order.DeliveryAddress = req.DeliveryAddress
	}
	if req.Notes != "" {
		order.Notes = req.Notes
	}

	order.LastModifiedBy = userID

	if err := database.DB.Save(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update order")
	}

	// Reload with associations
	database.DB.Preload("CreatedByUser").Preload("LastModifiedUser").First(&order, order.ID)

	return c.JSON(http.StatusOK, order)
}

// SoftDeleteOrder marks an order as deleted
func SoftDeleteOrder(c echo.Context) error {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	order.IsDeleted = true
	if err := database.DB.Save(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete order")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order deleted successfully"})
}

// RestoreOrder restores a soft-deleted order
func RestoreOrder(c echo.Context) error {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.Unscoped().Where("id = ? AND is_deleted = ?", id, true).First(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "deleted order not found")
	}

	order.IsDeleted = false
	if err := database.DB.Save(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to restore order")
	}

	return c.JSON(http.StatusOK, order)
}

// validateStatusTransition checks if a status transition is allowed for a given role
func validateStatusTransition(currentStatus, newStatus models.OrderStatus, role models.UserRole) error {
	transitions := map[models.OrderStatus]map[models.UserRole][]models.OrderStatus{
		models.StatusOrdered: {
			models.RoleWarehouse: {models.StatusInProcess},
		},
		models.StatusInProcess: {
			models.RoleWarehouse: {models.StatusInRoute},
		},
		models.StatusInRoute: {
			models.RoleRoute: {models.StatusDelivered},
		},
	}

	allowedStatuses, ok := transitions[currentStatus][role]
	if !ok {
		return fmt.Errorf("role %s cannot change status from %s", role, currentStatus)
	}

	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s for role %s", currentStatus, newStatus, role)
}
