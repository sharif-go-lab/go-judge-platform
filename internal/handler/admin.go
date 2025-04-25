// internal/handler/admin.go
package handler

import (
	"github.com/gin-contrib/sessions"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// AdminListUsersHandler lists all users for admin management
func AdminListUsersHandler(c *gin.Context) {
	var users []model.User
	if err := db.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Manage Users",
			"error": "Failed to fetch users",
		})
		return
	}

	userID := sessions.Default(c).Get("user_id")
	c.HTML(http.StatusOK, "users.html", gin.H{
		"userID":  userID,
		"isAdmin": true,
		"title":   "Manage Users",
		"users":   users,
	})
}

// PromoteUserHandler sets is_admin=true for a given user
func PromoteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.DB.Model(&model.User{}).
		Where("id = ?", id).
		Update("is_admin", true).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Manage Users",
			"error": "Failed to promote user",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// DemoteUserHandler sets is_admin=false for a given user
func DemoteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.DB.Model(&model.User{}).
		Where("id = ?", id).
		Update("is_admin", false).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Manage Users",
			"error": "Failed to demote user",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// PublishQuestionHandler sets is_published=true and published_at=now()
func PublishQuestionHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.DB.Model(&model.Problem{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"publish_date": time.Now(),
			"status":       "published",
		}).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to publish question",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/questions")
}

// UnpublishQuestionHandler sets is_published=false
func UnpublishQuestionHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.DB.Model(&model.Problem{}).
		Where("id = ?", id).
		Update("status", "unpublished").Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to unpublish question",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/questions")
}
