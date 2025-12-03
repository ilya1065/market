package migrations

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"marketplace/internal/config"
	"marketplace/internal/domain"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	conf, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке conf")
	}
	dsn := conf.DBDSN
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}

func Migration() {

	err := DB.AutoMigrate(domain.User{})
	if err != nil {
		log.Fatal(fmt.Errorf("не удалось осуществить миграцию %w", err))
	}
}
