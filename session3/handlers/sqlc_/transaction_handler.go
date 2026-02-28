package sqlc

import (
	"database/sql"
	"net/http"

	"gorm-sqlc/models"
	sqlcdb "gorm-sqlc/sqlc"

	"github.com/gin-gonic/gin"
)

type SQLCTransactionHandler struct {
	db      *sql.DB
	queries *sqlcdb.Queries
}

func NewSQLCTransactionHandler(db *sql.DB, queries *sqlcdb.Queries) *SQLCTransactionHandler {
	return &SQLCTransactionHandler{
		db:      db,
		queries: queries,
	}
}

// CreateUsersInTransaction handles POST /api/sqlc/transactions/users
func (h *SQLCTransactionHandler) CreateUsersInTransaction(c *gin.Context) {
	var req struct {
		Users []models.CreateUserRequest `json:"users" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)
	var createdUsers []sqlcdb.User

	for _, userReq := range req.Users {
		var age sql.NullInt32
		if userReq.Age != nil {
			age = sql.NullInt32{Int32: int32(*userReq.Age), Valid: true}
		}
		isActive := sql.NullBool{Bool: userReq.IsActive, Valid: true}
		user, err := qtx.CreateUser(ctx, sqlcdb.CreateUserParams{
			Email:    userReq.Email,
			Username: userReq.Username,
			FullName: userReq.FullName,
			Age:      age,
			IsActive: isActive,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		createdUsers = append(createdUsers, user)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "users created in transaction",
		"users":   createdUsers,
	})
}
