package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/middleware"
)

// RegisterProtectedRoutes sets up all routes that require a logged-in user
func RegisterProtectedRoutes(r *gin.Engine) {
	// Must apply AuthRequired BEFORE calling this in main.go
	r.Use(middleware.AuthRequired())

	// Question creation & “my questions”
	questions := r.Group("/questions")
	{
		questions.GET("/create", CreateQuestionPageHandler) // show create form
		questions.POST("/create", CreateQuestionHandler)    // submit new question
		questions.GET("/my", MyQuestionsHandler)            // list this user’s questions
		questions.GET("/edit/:id", EditQuestionPageHandler) // edit question
		questions.POST("/edit/:id", EditQuestionHandler)    // save edited question
	}

	// Submissions (history & submit new)
	submissions := r.Group("/submissions")
	{
		submissions.GET("/", MySubmissionsHandler)              // list this user’s submissions
		submissions.GET("/:id", SubmissionDetailHandler)        // view a single submission
		submissions.POST("/submit/:question_id", SubmitHandler) // submit code
	}

	r.GET("/profile", ProfileHandler)
}
