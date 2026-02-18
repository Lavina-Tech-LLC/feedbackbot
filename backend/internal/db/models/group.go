package models

import "gorm.io/gorm"

type (
	Group struct {
		TenantID uint   `gorm:"not null" json:"tenant_id"`
		BotID    uint   `gorm:"not null" json:"bot_id"`
		ChatID   int64  `gorm:"not null;uniqueIndex" json:"chat_id"`
		Title    string `json:"title"`
		Type     string `json:"type"` // group, supergroup, forum
		IsActive bool   `gorm:"default:true" json:"is_active"`
		Tenant   Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
		Bot      Bot    `gorm:"foreignKey:BotID" json:"bot,omitempty"`
		gorm.Model
	}

	FeedbackConfig struct {
		GroupID      uint  `gorm:"not null;uniqueIndex" json:"group_id"`
		PostToGroup  bool  `gorm:"default:false" json:"post_to_group"`
		ForumTopicID *int  `json:"forum_topic_id"`
		Group        Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
		gorm.Model
	}
)
