package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null;size:100" json:"name"`
	Email        string `gorm:"not null;unique;size:255" json:"email" validate:"required,email"`
	PasswordHash string `gorm:"not null" json:"-"`
	IsActive     bool   `gorm:"not null;default:false;index" json:"is_active"`
	IsAdmin      bool   `gorm:"not null;default:false" json:"is_admin"`

	Sections  []Section  `gorm:"foreignKey:UserID" json:"sections,omitempty"`
	Employees []Employee `gorm:"foreignKey:UserID" json:"employees,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginatedUsers struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	HasMore    bool           `json:"has_more"`
}
