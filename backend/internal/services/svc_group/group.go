package svc_group

import (
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/services"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	tenantID := services.GetTenantID(c)

	var groups []models.Group
	models.DB.Scopes(db.TenantScope(tenantID)).Find(&groups)

	c.Data(lvn.Res(200, groups, ""))
}

func GetGroup(c *gin.Context) {
	id := c.Param("id")

	var group models.Group
	if err := models.DB.First(&group, id).Error; err != nil {
		c.Data(lvn.Res(404, "", "Group not found"))
		return
	}

	c.Data(lvn.Res(200, group, ""))
}

type updateGroupReq struct {
	IsActive *bool `json:"is_active"`
}

func UpdateGroup(c *gin.Context) {
	id := c.Param("id")

	var group models.Group
	if err := models.DB.First(&group, id).Error; err != nil {
		c.Data(lvn.Res(404, "", "Group not found"))
		return
	}

	var req updateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "Invalid request: "+err.Error()))
		return
	}

	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}

	if err := models.DB.Save(&group).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to update group")
		return
	}

	c.Data(lvn.Res(200, group, ""))
}

type updateConfigReq struct {
	PostToGroup  *bool `json:"post_to_group"`
	ForumTopicID *int  `json:"forum_topic_id"`
}

func UpdateGroupConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.FeedbackConfig
	if err := models.DB.Where("group_id = ?", id).First(&config).Error; err != nil {
		c.Data(lvn.Res(404, "", "Config not found"))
		return
	}

	var req updateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Data(lvn.Res(400, "", "Invalid request: "+err.Error()))
		return
	}

	if req.PostToGroup != nil {
		config.PostToGroup = *req.PostToGroup
	}
	if req.ForumTopicID != nil {
		config.ForumTopicID = req.ForumTopicID
	}

	if err := models.DB.Save(&config).Error; err != nil {
		lvn.GinErr(c, 500, err, "Failed to update config")
		return
	}

	c.Data(lvn.Res(200, config, ""))
}
