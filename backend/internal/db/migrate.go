package db

import "github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"

func Migrate() {
	err := models.DB.AutoMigrate(
		&models.Tenant{},
		&models.Bot{},
		&models.Group{},
		&models.FeedbackConfig{},
	)
	if err != nil {
		panic(err)
	}
}
