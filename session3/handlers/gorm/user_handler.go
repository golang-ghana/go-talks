package gorm

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm-sqlc/models"
)

type GormUserHandler struct {
	db *gorm.DB
}

func NewGormUserHandler(db *gorm.DB) *GormUserHandler {
	return &GormUserHandler{db: db}
}

// CreateUser handles POST /api/gorm/users
func (h *GormUserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Age:      req.Age,
		IsActive: req.IsActive,
	}

	result := h.db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// GetUserByID handles GET /api/gorm/users/:id
func (h *GormUserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user models.User
	result := h.db.First(&user, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByEmail handles GET /api/gorm/users/email/:email
func (h *GormUserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	var user models.User
	result := h.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers handles GET /api/gorm/users
func (h *GormUserHandler) ListUsers(c *gin.Context) {
	var req models.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Limit = 10
		req.Offset = 0
	}

	var users []models.User
	result := h.db.Order("created_at DESC").Limit(req.Limit).Offset(req.Offset).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ListActiveUsers handles GET /api/gorm/users/active
func (h *GormUserHandler) ListActiveUsers(c *gin.Context) {
	var users []models.User
	result := h.db.Where("is_active = ?", true).Order("created_at DESC").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles PUT /api/gorm/users/:id
// Note: Uses Save() which requires SELECT first, then UPDATE (2 queries)
func (h *GormUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Email = req.Email
	user.Username = req.Username
	user.FullName = req.FullName
	user.Age = req.Age

	result := h.db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserStatus handles PATCH /api/gorm/users/:id/status
func (h *GormUserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", uint(id)).Updates(map[string]interface{}{
		"is_active":  req.IsActive,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}

// DeleteUser handles DELETE /api/gorm/users/:id
func (h *GormUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	result := h.db.Delete(&models.User{}, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// SearchUsersByName handles GET /api/gorm/users/search
func (h *GormUserHandler) SearchUsersByName(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	var users []models.User
	result := h.db.Where("full_name ILIKE ?", "%"+req.Query+"%").Order("full_name").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUsersByAgeRange handles GET /api/gorm/users/age-range
func (h *GormUserHandler) GetUsersByAgeRange(c *gin.Context) {
	var req models.AgeRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []models.User
	result := h.db.Where("age BETWEEN ? AND ?", req.MinAge, req.MaxAge).Order("age").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CountUsers handles GET /api/gorm/users/count
func (h *GormUserHandler) CountUsers(c *gin.Context) {
	var count int64
	result := h.db.Model(&models.User{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
