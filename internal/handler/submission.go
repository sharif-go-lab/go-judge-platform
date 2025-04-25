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
	userID := sessions.Default(c).Get("user_id")

	// Query submissions for this user
	var submissions []model.Submission
	db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&submissions)

	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	c.HTML(http.StatusOK, "sList.html", gin.H{
		"userID":      userID,
		"isAdmin":     isAdmin,
		"title":       "My Submissions",
		"submissions": submissions,
	})
}

// SubmissionDetailHandler displays the details of a single submission
func SubmissionDetailHandler(c *gin.Context) {
	// Parse submission ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Submission not found",
		})
		return
	}

	// Fetch submission record
	var submission model.Submission
	if err := db.DB.First(&submission, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Submission not found",
		})
		return
	}

	// Ensure only the owner or an admin can view
	userID := sessions.Default(c).Get("user_id")
	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	if submission.UserID != userID && !isAdmin {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Access denied",
		})
		return
	}

	c.HTML(http.StatusOK, "sDetail.html", gin.H{
		"userID":     userID,
		"isAdmin":    isAdmin,
		"title":      "Submission Details",
		"submission": submission,
	})
}

// SubmitHandler processes code submissions for a question
func SubmitHandler(c *gin.Context) {
	// Parse question ID
	questionID, err := strconv.Atoi(c.Param("question_id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Question not found",
		})
		return
	}

	// Bind submission form (code text)
	var form struct {
		Code string `form:"code" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "My Submissions",
			"error": err.Error(),
		})
		return
	}

	// Create submission record
	userID := sessions.Default(c).Get("user_id").(uint)
	sub := model.Submission{
		UserID:    userID,
		ProblemID: uint(questionID),
		Code:      form.Code,
		Language:  "Golang",
		Status:    model.StatusPending,
	}
	if err := db.DB.Create(&sub).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Failed to submit code",
		})
		return
	}

	// Redirect to the submission detail page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/submissions/%d", sub.ID))
}
