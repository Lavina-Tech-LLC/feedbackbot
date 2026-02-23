package db

import (
	"fmt"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() {
	cfg := config.Confs.DB
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Tashkent",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	models.DB = db
	Migrate()
}
