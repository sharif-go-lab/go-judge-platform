package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes sets up all routes that don't require authentication
func RegisterPublicRoutes(r *gin.Engine) {
	// Home page
	r.GET("/", HomeHandler)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.GET("/login", LoginPageHandler)
		auth.POST("/login", LoginHandler)
		auth.GET("/register", RegisterPageHandler)
		auth.POST("/register", RegisterHandler)
		auth.GET("/logout", LogoutHandler)
	}

	// Public question listing & detail
	questions := r.Group("/questions")
	{
		questions.GET("/", QuestionListHandler)      // list all published questions
		questions.GET("/:id", QuestionDetailHandler) // view a single question
	}

	// **Profile view** is public (so only registered once)
	r.GET("/profile/:username", ProfileHandler)
}
