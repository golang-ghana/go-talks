package gorm

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm-sqlc/models"
)

type GormAdvancedHandler struct {
	db *gorm.DB
}

func NewGormAdvancedHandler(db *gorm.DB) *GormAdvancedHandler {
	return &GormAdvancedHandler{db: db}
}

// ComplexQuery handles GET /api/gorm/users/complex
func (h *GormAdvancedHandler) ComplexQuery(c *gin.Context) {
	var users []models.User
	result := h.db.Where("is_active = ?", true).
		Where("age >= ?", 25).
		Or("username LIKE ?", "admin%").
		Order("created_at DESC").
		Limit(5).
		Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// SelectSpecificFields handles GET /api/gorm/users/partial
func (h *GormAdvancedHandler) SelectSpecificFields(c *gin.Context) {
	var partialUsers []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	result := h.db.Model(&models.User{}).Select("id", "username", "email").Find(&partialUsers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, partialUsers)
}

// RawSQLQuery handles GET /api/gorm/users/raw
func (h *GormAdvancedHandler) RawSQLQuery(c *gin.Context) {
	age := c.DefaultQuery("min_age", "20")
	limit := c.DefaultQuery("limit", "5")

	var users []models.User
	result := h.db.Raw("SELECT * FROM users WHERE age > ? ORDER BY age LIMIT ?", age, limit).Scan(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
