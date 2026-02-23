package svc_feedback

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

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

func applyFilters(c *gin.Context) *gorm.DB {
	tenantID := services.GetTenantID(c)
	query := models.DB.Scopes(db.TenantScope(tenantID))

	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	if adminOnly := c.Query("admin_only"); adminOnly == "true" {
		query = query.Where("admin_only = ?", true)
	} else if adminOnly == "false" {
		query = query.Where("admin_only = ?", false)
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("created_at <= ?", dateTo)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("message LIKE ?", "%"+search+"%")
	}

	return query
}

func GetFeedbacks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := applyFilters(c)

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

func ExportCSV(c *gin.Context) {
	query := applyFilters(c)

	var feedbacks []models.Feedback
	query.Preload("Group", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "title")
	}).Order("created_at DESC").Limit(10000).Find(&feedbacks)

	filename := fmt.Sprintf("feedbacks_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	w := csv.NewWriter(c.Writer)
	w.Write([]string{"ID", "Group", "Message", "Admin Only", "Posted to Group", "Created At"})

	for _, fb := range feedbacks {
		w.Write([]string{
			strconv.FormatUint(uint64(fb.ID), 10),
			fb.Group.Title,
			fb.Message,
			strconv.FormatBool(fb.AdminOnly),
			strconv.FormatBool(fb.Posted),
			fb.CreatedAt.Format(time.RFC3339),
		})
	}

	w.Flush()
}
