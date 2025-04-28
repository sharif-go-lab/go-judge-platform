package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// ProfileHandler displays a user's profile pulled from the database
func ProfileHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		usernameStr := sessions.Default(c).Get("username")
		if usernameStr == nil {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"title":   "Unauthorized User",
				"message": "You should login to see your profile",
			})
			return
		}
		username = usernameStr.(string)
	}

	// Fetch user by username
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		// user not found: show 404 page
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "User Not Found",
			"message": "The requested user does not exist.",
		})
		return
	}

	// Compute stats: attempted and solved submissions
	var attempted int64
	db.DB.Model(&model.Submission{}).
		Where("user_id = ?", user.ID).
		Count(&attempted)

	var solved int64
	db.DB.Model(&model.Submission{}).
		Where("user_id = ? AND status = ?", user.ID, "accepted").
		Count(&solved)

	// Calculate success rate
	var successRate float64
	if attempted > 0 {
		successRate = (float64(solved) / float64(attempted)) * 100
	}

	isAdmin, userID := false, sessions.Default(c).Get("user_id")
	if sessions.Default(c).Get("is_admin") != nil {
		isAdmin = sessions.Default(c).Get("is_admin").(bool)
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"userID":  userID,
		"isAdmin": isAdmin,
		"title":   user.Username + "'s Profile",
		"user":    user,
		"stats": gin.H{
			"attempted":   attempted,
			"solved":      solved,
			"successRate": fmt.Sprintf("%.2f", successRate),
		},
	})
}
