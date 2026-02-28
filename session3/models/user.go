package models

import (
	"time"
)

// User model for GORM
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username  string    `gorm:"uniqueIndex;not null;size:50" json:"username"`
	FullName  string    `gorm:"not null;size:100" json:"full_name"`
	Age       *int      `gorm:"default:null" json:"age,omitempty"`
	IsActive  bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Request/Response DTOs
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Age      *int   `json:"age,omitempty"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Age      *int   `json:"age,omitempty"`
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"`
}

type ListUsersRequest struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

type AgeRangeRequest struct {
	MinAge int `form:"min_age" binding:"required,min=0"`
	MaxAge int `form:"max_age" binding:"required,min=0"`
}

type SearchRequest struct {
	Query string `form:"q" binding:"required"`
}
