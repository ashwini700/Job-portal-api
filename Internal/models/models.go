package models

import (
	"gorm.io/gorm"
)

// user is an employee
type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type NewUser struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Company struct {
	gorm.Model
	// CompanyId   uint   `json:"companyId"`
	CompanyName string `json:"companyname"`
	Location    string `json:"location"`
	Jobs []Job  `json:"jobs,omitempty" gorm:"foreignKey:CompanyId"`
}

type Job struct {
	gorm.Model
	Name       string `json:"title"`
	Field      string `json:"field"`
	Experience uint   `json:"experience"`
	CompanyId  uint64 `json:"companyId"`
}
