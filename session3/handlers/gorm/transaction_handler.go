package gorm

import (
	"net/http"

	"gorm-sqlc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GormTransactionHandler struct {
	db *gorm.DB
}

func NewGormTransactionHandler(db *gorm.DB) *GormTransactionHandler {
	return &GormTransactionHandler{db: db}
}

// CreateUsersInTransaction handles POST /api/gorm/transactions/users
func (h *GormTransactionHandler) CreateUsersInTransaction(c *gin.Context) {
	var req struct {
		Users []models.CreateUserRequest `json:"users" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createdUsers []models.User
	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, userReq := range req.Users {
			user := models.User{
				Email:    userReq.Email,
				Username: userReq.Username,
				FullName: userReq.FullName,
				Age:      userReq.Age,
				IsActive: userReq.IsActive,
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			createdUsers = append(createdUsers, user)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "users created in transaction",
		"users":   createdUsers,
	})
}
