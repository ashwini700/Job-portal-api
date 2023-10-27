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
	CompanyName string `json:"companyname"`
	JobRole     string `json:"jobRole"`
}

// employee contains information needed to create a Employee details.
// type Employee struct {
// 	EmpName string  `json:"emp_name" validate:"required"`
// 	EmpRole int     `json:"emprole" validate:"required,number"`
// 	EmpId   float64 `json:"emp_id" validate:"required,number,gt=0"`
// }
