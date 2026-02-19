package svc_feedback

import (
	"strconv"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedbackResponse struct {
	models.Feedback
	GroupName string `json:"group_name"`
}

func GetFeedbacks(c *gin.Context) {
	adminOnly := c.Query("admin_only")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	tenantID := services.GetTenantID(c)
	query := models.DB.Scopes(db.TenantScope(tenantID))

	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	if adminOnly == "true" {
		query = query.Where("admin_only = ?", true)
	} else if adminOnly == "false" {
		query = query.Where("admin_only = ?", false)
	}

	if dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_at <= ?", dateTo)
	}

	var total int64
	query.Model(&models.Feedback{}).Count(&total)

	var feedbacks []models.Feedback
	query.Preload("Group", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "title")
	}).Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks)

	resp := make([]FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		resp[i] = FeedbackResponse{Feedback: fb, GroupName: fb.Group.Title}
	}

	c.Data(lvn.Res(200, gin.H{
		"data":  resp,
		"total": total,
		"page":  page,
		"limit": limit,
	}, ""))
}
