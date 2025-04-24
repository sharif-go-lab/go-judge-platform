package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// QuestionListHandler displays the list of published questions with pagination.
func QuestionListHandler(c *gin.Context) {
	// parse page number
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	const perPage = 10

	// count total published problems
	var total int64
	db.DB.Model(&model.Problem{}).
		Where("status = ?", "published").
		Count(&total)

	totalPages := int((total + perPage - 1) / perPage)

	// fetch current page
	var questions []model.Problem
	db.DB.
		Where("status = ?", "published").
		Order("publish_date DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&questions)

	c.HTML(http.StatusOK, "qList.html", gin.H{
		"title":      "Questions",
		"questions":  questions,
		"page":       page,
		"totalPages": totalPages,
		"prevPage":   page - 1,
		"nextPage":   page + 1,
	})
}

// QuestionDetailHandler displays a single published question.
func QuestionDetailHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Question not found"})
		return
	}

	var question model.Problem
	if err := db.DB.First(&question, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Question not found"})
		return
	}
	if question.Status != "published" {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Question not available"})
		return
	}

	c.HTML(http.StatusOK, "qDetail.html", gin.H{
		"title":    question.Title,
		"question": question,
	})
}

// CreateQuestionPageHandler shows the new-question form.
func CreateQuestionPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{"title": "Create Question"})
}

// CreateQuestionHandler processes the new-question form.
func CreateQuestionHandler(c *gin.Context) {
	var form struct {
		Title         string `form:"title" binding:"required"`
		Statement     string `form:"statement" binding:"required"`
		TimeLimitMs   int    `form:"time_limit_ms" binding:"required"`
		MemoryLimitMb int    `form:"memory_limit_mb" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"title": "Create Question",
			"error": err.Error(),
		})
		return
	}

	user := c.MustGet("user").(model.User)
	question := model.Problem{
		OwnerID:       user.ID,
		Title:         form.Title,
		Statement:     form.Statement,
		TimeLimitMs:   form.TimeLimitMs,
		MemoryLimitMb: form.MemoryLimitMb,
		Status:        "draft",
	}

	if err := db.DB.Create(&question).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "create.html", gin.H{
			"title": "Create Question",
			"error": "Failed to save question",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/questions")
}

// MyQuestionsHandler shows all questions the current user created.
func MyQuestionsHandler(c *gin.Context) {
	user := c.MustGet("user").(model.User)

	var questions []model.Problem
	db.DB.
		Where("owner_id = ?", user.ID).
		Order("created_at DESC").
		Find(&questions)

	c.HTML(http.StatusOK, "my-questions.html", gin.H{
		"title":     "My Questions",
		"questions": questions,
	})
}
