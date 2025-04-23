package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// ProfileHandler displays a user's profile
func ProfileHandler(c *gin.Context) {
	username := c.Param("username")

	// For Phase 1, generate a mock user profile
	user := model.User{
		ID:        1,
		Username:  username,
		Email:     username + "@example.com",
		IsAdmin:   false,
		CreatedAt: time.Now().AddDate(0, -3, 0), // 3 months ago
	}

	// Mock stats
	stats := gin.H{
		"attempted":   15,
		"solved":      10,
		"successRate": 66.67, // 10/15 * 100
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"title": username + "'s Profile",
		"user":  user,
		"stats": stats,
	})
}
