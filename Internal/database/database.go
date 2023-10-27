package database

import (
	"Job-portal-api/Internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Jobs struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique;not null"`
	// Role    string
	// Company string
}

var db *gorm.DB

func Open() (*gorm.DB, error) {
	dataSources := "host=localhost user=postgres password=Ashwini dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dataSources), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Migrator().AutoMigrate(&models.User{}, &models.Company{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
