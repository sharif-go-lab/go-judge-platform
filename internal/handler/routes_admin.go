// internal/handler/routes_admin.go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/middleware"
)

// RegisterAdminRoutes mounts admin-only endpoints on the given Gin engine.
// It applies authentication and admin-role checks internally.
func RegisterAdminRoutes(r *gin.Engine) {
	// Admin group under /admin
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())

	// User management
	admin.GET("/users", AdminListUsersHandler)
	admin.POST("/users/promote/:id", PromoteUserHandler)
	admin.POST("/users/demote/:id", DemoteUserHandler)

	// Question management
	admin.POST("/questions/publish/:id", PublishQuestionHandler)
	admin.POST("/questions/unpublish/:id", UnpublishQuestionHandler)
}
