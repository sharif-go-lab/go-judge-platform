// internal/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// AuthRequired ensures the user is logged in
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		uid := sess.Get("user_id")
		if uid == nil {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}
		var user model.User
		if err := db.DB.First(&user, uid).Error; err != nil {
			sess.Clear()
			sess.Save()
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}
		c.Set("currentUser", &user)
		c.Next()
	}
}

// AdminRequired ensures the user is an admin
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cu, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		user := cu.(*model.User)
		if !user.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
