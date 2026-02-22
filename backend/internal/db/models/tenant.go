package models

import "gorm.io/gorm"

type (
	Tenant struct {
		Name string `gorm:"not null" json:"name"`
		Slug string `gorm:"uniqueIndex;not null" json:"slug"`
		Bots []Bot  `gorm:"foreignKey:TenantID" json:"bots,omitempty"`
		gorm.Model
	}

	Bot struct {
		TenantID    uint   `gorm:"not null" json:"tenant_id"`
		Token       string `gorm:"not null" json:"-"`
		BotUsername string `json:"bot_username"`
		BotName     string `json:"bot_name"`
		Verified    bool   `gorm:"default:false" json:"verified"`
		Tenant      Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
		gorm.Model
	}

	User struct {
		Email        string `gorm:"uniqueIndex;not null" json:"email"`
		Name         string `gorm:"not null" json:"name"`
		PasswordHash string `gorm:"not null" json:"-"`
		Role         string `gorm:"default:user" json:"role"`
		gorm.Model
	}

	UserTenant struct {
		UserID   string `gorm:"uniqueIndex;not null" json:"user_id"`
		TenantID uint   `gorm:"not null" json:"tenant_id"`
		Tenant   Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
		gorm.Model
	}
)
