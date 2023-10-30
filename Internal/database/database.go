package database

import (
	"Job-portal-api/Internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var db *gorm.DB

func Open() (*gorm.DB, error) {
	dataSources := "host=localhost user=postgres password=Ashwini dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dataSources), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Migrator().AutoMigrate(&models.User{}, &models.Company{}, &models.Job{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
