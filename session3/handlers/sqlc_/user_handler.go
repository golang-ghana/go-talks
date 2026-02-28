package sqlc

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"gorm-sqlc/models"
	sqlcdb "gorm-sqlc/sqlc"

	"github.com/gin-gonic/gin"
)

// Querier interface for SQLC queries
type Querier interface {
	CreateUser(ctx context.Context, arg sqlcdb.CreateUserParams) (sqlcdb.User, error)
	GetUserByID(ctx context.Context, id int32) (sqlcdb.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlcdb.User, error)
	ListUsers(ctx context.Context, arg sqlcdb.ListUsersParams) ([]sqlcdb.User, error)
	ListActiveUsers(ctx context.Context) ([]sqlcdb.User, error)
	UpdateUser(ctx context.Context, arg sqlcdb.UpdateUserParams) (sqlcdb.User, error)
	UpdateUserActiveStatus(ctx context.Context, arg sqlcdb.UpdateUserActiveStatusParams) error
	DeleteUser(ctx context.Context, id int32) error
	CountUsers(ctx context.Context) (int64, error)
	SearchUsersByName(ctx context.Context, query sql.NullString) ([]sqlcdb.User, error)
	GetUsersByAgeRange(ctx context.Context, arg sqlcdb.GetUsersByAgeRangeParams) ([]sqlcdb.User, error)
}

type SQLCUserHandler struct {
	queries Querier
}

func NewSQLCUserHandler(queries Querier) *SQLCUserHandler {
	return &SQLCUserHandler{
		queries: queries,
	}
}

// CreateUser handles POST /api/sqlc/users
func (h *SQLCUserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var age sql.NullInt32
	if req.Age != nil {
		age = sql.NullInt32{Int32: int32(*req.Age), Valid: true}
	}
	isActive := sql.NullBool{Bool: req.IsActive, Valid: true}
	newUser, err := h.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Age:      age,
		IsActive: isActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

// GetUserByID handles GET /api/sqlc/users/:id
func (h *SQLCUserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByEmail handles GET /api/sqlc/users/email/:email
func (h *SQLCUserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	ctx := c.Request.Context()
	user, err := h.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUsers handles GET /api/sqlc/users
func (h *SQLCUserHandler) ListUsers(c *gin.Context) {
	var req models.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Limit = 10
		req.Offset = 0
	}

	ctx := c.Request.Context()
	users, err := h.queries.ListUsers(ctx, sqlcdb.ListUsersParams{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// ListActiveUsers handles GET /api/sqlc/users/active
func (h *SQLCUserHandler) ListActiveUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.queries.ListActiveUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser handles PUT /api/sqlc/users/:id
func (h *SQLCUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var age sql.NullInt32
	if req.Age != nil {
		age = sql.NullInt32{Int32: int32(*req.Age), Valid: true}
	}
	updatedUser, err := h.queries.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:       int32(id),
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Age:      age,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// UpdateUserStatus handles PATCH /api/sqlc/users/:id/status
func (h *SQLCUserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	isActive := sql.NullBool{Bool: req.IsActive, Valid: true}
	err = h.queries.UpdateUserActiveStatus(ctx, sqlcdb.UpdateUserActiveStatusParams{
		ID:       int32(id),
		IsActive: isActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}

// DeleteUser handles DELETE /api/sqlc/users/:id
func (h *SQLCUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.queries.DeleteUser(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// SearchUsersByName handles GET /api/sqlc/users/search
func (h *SQLCUserHandler) SearchUsersByName(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	ctx := c.Request.Context()
	query := sql.NullString{String: req.Query, Valid: true}
	users, err := h.queries.SearchUsersByName(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUsersByAgeRange handles GET /api/sqlc/users/age-range
func (h *SQLCUserHandler) GetUsersByAgeRange(c *gin.Context) {
	var req models.AgeRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	users, err := h.queries.GetUsersByAgeRange(ctx, sqlcdb.GetUsersByAgeRangeParams{
		Age:   sql.NullInt32{Int32: int32(req.MinAge), Valid: true},
		Age_2: sql.NullInt32{Int32: int32(req.MaxAge), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CountUsers handles GET /api/sqlc/users/count
func (h *SQLCUserHandler) CountUsers(c *gin.Context) {
	ctx := c.Request.Context()
	count, err := h.queries.CountUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
