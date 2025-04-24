package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// MySubmissionsHandler displays the current user’s submission history
func MySubmissionsHandler(c *gin.Context) {
	// Get current user from context (set by auth middleware)
	sess := sessions.Default(c)
	userID := sess.Get("user_ID")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}

	// Query submissions for this user
	var submissions []model.Submission
	db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&submissions)

	c.HTML(http.StatusOK, "sList.html", gin.H{
		"title":       "My Submissions",
		"submissions": submissions,
	})
}

// SubmissionDetailHandler displays the details of a single submission
func SubmissionDetailHandler(c *gin.Context) {
	// Parse submission ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Submission not found"})
		return
	}

	// Fetch submission record
	var submission model.Submission
	if err := db.DB.First(&submission, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Submission not found"})
		return
	}

	// Ensure only the owner or an admin can view
	user := c.MustGet("user").(model.User)
	if submission.UserID != user.ID && !user.IsAdmin {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Access denied"})
		return
	}

	c.HTML(http.StatusOK, "submission/sDetail.html", gin.H{
		"title":      "Submission Details",
		"submission": submission,
	})
}

// SubmitHandler processes code submissions for a question
func SubmitHandler(c *gin.Context) {
	// Parse question ID
	questionID, err := strconv.Atoi(c.Param("question_id"))
	if err != nil || questionID < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Question not found"})
		return
	}

	// Bind submission form (code text + language)
	var form struct {
		Code     string `form:"code" binding:"required"`
		Language string `form:"language" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "submission/submit.html", gin.H{
			"error":    "Code and language are required",
			"question": questionID,
		})
		return
	}

	// Create submission record
	user := c.MustGet("user").(model.User)
	sub := model.Submission{
		UserID:    user.ID,
		ProblemID: uint(questionID),
		Code:      form.Code,
		Language:  form.Language,
		Status:    model.StatusPending,
	}
	if err := db.DB.Create(&sub).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "submission/submit.html", gin.H{
			"error":    "Failed to submit code",
			"question": questionID,
		})
		return
	}

	// Redirect to the submission detail page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/submissions/%d", sub.ID))
}
