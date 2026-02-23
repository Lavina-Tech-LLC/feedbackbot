package db

import "github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"

func Migrate() {
	err := models.DB.AutoMigrate(
		&models.Tenant{},
		&models.Bot{},
		&models.Group{},
		&models.FeedbackConfig{},
		&models.GroupUser{},
		&models.Feedback{},
		&models.User{},
		&models.UserTenant{},
		&models.PendingFeedback{},
	)
	if err != nil {
		panic(err)
	}
}
