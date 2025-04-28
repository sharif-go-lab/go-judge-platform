package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/spf13/viper"

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

	var question model.Problem
	if err := db.DB.First(&question, questionID).Error; err != nil {
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
		Status:    "Pending",
	}
	if err := db.DB.Create(&sub).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "My Submissions",
			"error": "Failed to submit code",
		})
		return
	}

	go func() {
		runnerURL := fmt.Sprintf("http://code-runner%s/run", viper.GetString("code_runner.listen"))
		req := model.RunRequest{
			Code:          sub.Code,
			SampleInput:   question.SampleInput,
			SampleOutput:  question.SampleOutput,
			TimeLimitMs:   question.TimeLimitMs,
			MemoryLimitMb: question.MemoryLimitMb,
		}
		jsonData, err := json.Marshal(req)
		if err != nil {
			log.Println("Failed to marshal submission", sub.ID, "data", err)
			if err := db.DB.Model(&sub).Update("status", "Failed").Error; err != nil {
				log.Println("Failed to update submission", sub.ID, "status", "Failed", "error", err)
			}
			return
		}

		resp, err := http.Post(runnerURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Failed to send submission", sub.ID, "request", err)
			if err := db.DB.Model(&sub).Update("status", "Failed").Error; err != nil {
				log.Println("Failed to update submission", sub.ID, "status", "Failed", "error", err)
			}
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read submission", sub.ID, "response", err)
			if err := db.DB.Model(&sub).Update("status", "Failed").Error; err != nil {
				log.Println("Failed to update submission", sub.ID, "status", "Failed", "error", err)
			}
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Println("Failed to unmarshal submission", sub.ID, "response", err)
			if err := db.DB.Model(&sub).Update("status", "Failed").Error; err != nil {
				log.Println("Failed to update submission", sub.ID, "status", "Failed", "error", err)
			}
			return
		}

		if err := db.DB.Model(&sub).Update("status", result["result"]).Error; err != nil {
			log.Println("Failed to update submission", sub.ID, "status", result["result"], "error", err)
			return
		}
		log.Println("Submission", sub.ID, "status updated", result["result"])
	}()
	go func() {
		timeOut := viper.GetInt("server.submission_time_out")
		fmt.Println("Time out", timeOut)
		time.Sleep(time.Duration(timeOut) * time.Second)
		if err := db.DB.First(&sub, sub.ID).Error; err != nil {
			log.Println("Failed to get submission", sub.ID, "error", err)
			return
		}
		if sub.Status == "Pending" {
			if err := db.DB.Model(&sub).Update("status", "Failed").Error; err != nil {
				log.Println("Failed to update submission", sub.ID, "status", "Failed", "error", err)
			}
			log.Println("Submission", sub.ID, "status updated automatically to Failed")
		}
	}()

	// Redirect to the submission detail page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/submissions/%d", sub.ID))
}
