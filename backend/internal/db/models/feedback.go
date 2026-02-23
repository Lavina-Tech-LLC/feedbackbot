package models

import "gorm.io/gorm"

type (
	GroupUser struct {
		TenantID       uint  `gorm:"not null" json:"tenant_id"`
		GroupID        uint  `gorm:"not null" json:"group_id"`
		TelegramUserID int64 `gorm:"not null" json:"telegram_user_id"`
		Group          Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
		gorm.Model
	}

	Feedback struct {
		TenantID  uint      `gorm:"not null" json:"tenant_id"`
		GroupID   uint      `gorm:"not null" json:"group_id"`
		SenderID  uint      `gorm:"not null" json:"-"` // Never exposed via API
		Message   string    `gorm:"not null" json:"message"`
		AdminOnly bool      `gorm:"default:false" json:"admin_only"`
		Posted    bool      `gorm:"default:false" json:"posted"`
		Group     Group     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
		Sender    GroupUser `gorm:"foreignKey:SenderID" json:"-"` // Never exposed
		gorm.Model
	}

	PendingFeedback struct {
		TelegramUserID int64  `gorm:"uniqueIndex;not null" json:"telegram_user_id"`
		BotID          uint   `gorm:"not null" json:"bot_id"`
		Text           string `gorm:"not null" json:"text"`
		AdminOnly      bool   `gorm:"default:false" json:"admin_only"`
		gorm.Model
	}
)
