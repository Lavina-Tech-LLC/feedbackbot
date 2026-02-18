package svc_feedback

import (
	"strconv"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetFeedbacks(c *gin.Context) {
	groupID := c.Query("group_id")
	if groupID == "" {
		c.Data(lvn.Res(400, "", "group_id is required"))
		return
	}

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

	query := models.DB.Where("group_id = ?", groupID)

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
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks)

	c.Data(lvn.Res(200, gin.H{
		"data":  feedbacks,
		"total": total,
		"page":  page,
		"limit": limit,
	}, ""))
}
