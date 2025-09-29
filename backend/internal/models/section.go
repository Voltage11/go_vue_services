package models

import "gorm.io/gorm"

type Section struct {
	gorm.Model
	Name    string `gorm:"not null;size:255;index" json:"name"`
	Comment string `gorm:"type:text" json:"comment"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Связь многие-ко-многим с сотрудниками
	Employees []Employee `gorm:"many2many:employee_sections;" json:"employees,omitempty"`
}

func (s *Section) TableName() string {
	return "sections"
}
