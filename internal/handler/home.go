package handler

import (
	"github.com/gin-contrib/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler handles the home page
func HomeHandler(c *gin.Context) {
	userID := sessions.Default(c).Get("user_id")
	if userID != nil {
		c.Redirect(http.StatusSeeOther, "/profile")
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Online Judge System",
	})
}
