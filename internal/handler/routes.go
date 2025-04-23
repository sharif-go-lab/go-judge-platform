package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all the routes for the application
func RegisterRoutes(r *gin.Engine) {
	// Public routes
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

	// Questions routes
	questions := r.Group("/questions")
	{
		questions.GET("/", QuestionListHandler)
		questions.GET("/:id", QuestionDetailHandler)
		questions.GET("/create", CreateQuestionPageHandler)
		questions.POST("/create", CreateQuestionHandler)
		questions.GET("/my", MyQuestionsHandler)
	}

	// Submissions routes
	submissions := r.Group("/submissions")
	{
		submissions.GET("/", MySubmissionsHandler)
		submissions.GET("/:id", SubmissionDetailHandler)
		submissions.POST("/submit/:question_id", SubmitHandler)
	}

	// Profile routes
	r.GET("/profile/:username", ProfileHandler)
}
