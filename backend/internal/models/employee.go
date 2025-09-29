package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name     string `gorm:"not null;size:255;index" json:"name"`
	Email    string `gorm:"size:255;index" json:"email" validate:"omitempty,email"`
	Phone    string `gorm:"size:20" json:"phone"`
	IsActive bool   `gorm:"not null;default:false;index" json:"is_active"`

	// Принадлежит пользователю (владелец)
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Связь многие-ко-многим с секциями
	Sections []Section `gorm:"many2many:employee_sections;" json:"sections,omitempty"`
}
