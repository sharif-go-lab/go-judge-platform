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

// QuestionListHandler displays the list of published questions with pagination.
func QuestionListHandler(c *gin.Context) {
	// parse page number
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	const perPage = 10

	isAdmin := sessions.Default(c).Get("is_admin").(bool)

	// count total published problems
	var total int64
	if isAdmin == true {
		db.DB.Model(&model.Problem{}).Count(&total)
	} else {
		db.DB.Model(&model.Problem{}).
			Where("status = ?", "published").
			Count(&total)
	}

	totalPages := int((total + perPage - 1) / perPage)

	// fetch current page
	var questions []model.Problem
	if isAdmin == true {
		db.DB.Order("publish_date DESC").
			Offset((page - 1) * perPage).
			Limit(perPage).
			Find(&questions)
	} else {
		db.DB.Where("status = ?", "published").
			Order("publish_date DESC").
			Offset((page - 1) * perPage).
			Limit(perPage).
			Find(&questions)
	}

	userID := sessions.Default(c).Get("user_id")
	c.HTML(http.StatusOK, "qList.html", gin.H{
		"userID":     userID,
		"isAdmin":    isAdmin,
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
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Questions",
			"error": "Question not found",
		})
		return
	}

	var question model.Problem
	userID := sessions.Default(c).Get("user_id")
	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	if err := db.DB.First(&question, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Questions",
			"error": "Question not found",
		})
		return
	}
	if question.Status != "published" && !isAdmin && question.OwnerID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title": "Questions",
			"error": "Question not available",
		})
		return
	}

	c.HTML(http.StatusOK, "qDetail.html", gin.H{
		"userID":   userID,
		"isAdmin":  isAdmin,
		"title":    question.Title,
		"question": question,
	})
}

// CreateQuestionPageHandler shows the new-question form.
func CreateQuestionPageHandler(c *gin.Context) {
	userID := sessions.Default(c).Get("user_id")
	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	c.HTML(http.StatusOK, "create.html", gin.H{
		"userID":  userID,
		"isAdmin": isAdmin,
		"title":   "Create Question",
	})
}

// CreateQuestionHandler processes the new-question form.
func CreateQuestionHandler(c *gin.Context) {
	var form struct {
		Title         string `form:"title" binding:"required"`
		Statement     string `form:"statement" binding:"required"`
		TimeLimitMs   int    `form:"time_limit" binding:"required"`
		MemoryLimitMb int    `form:"memory_limit" binding:"required"`
		SampleInput   string `form:"input_test" binding:"required"`
		SampleOutput  string `form:"expected_output" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Create Question",
			"error": err.Error(),
		})
		return
	}

	userID := sessions.Default(c).Get("user_id").(uint)
	question := model.Problem{
		OwnerID:       userID,
		Title:         form.Title,
		Statement:     form.Statement,
		TimeLimitMs:   form.TimeLimitMs,
		MemoryLimitMb: form.MemoryLimitMb,
		Status:        "draft",
		SampleInput:   form.SampleInput,
		SampleOutput:  form.SampleOutput,
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

func EditQuestionPageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not found",
		})
		return
	}

	var question model.Problem
	userID := sessions.Default(c).Get("user_id")
	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	if err := db.DB.First(&question, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not found",
		})
		return
	}

	if question.OwnerID != userID && !isAdmin {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not available",
		})
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"userID":   userID,
		"isAdmin":  isAdmin,
		"title":    "Edit Question",
		"question": question,
	})
}

func EditQuestionHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not found",
		})
	}

	var form struct {
		Title         string `form:"title" binding:"required"`
		Statement     string `form:"statement" binding:"required"`
		TimeLimitMs   int    `form:"time_limit" binding:"required"`
		MemoryLimitMs int    `form:"memory_limit" binding:"required"`
		SampleInput   string `form:"input_test" binding:"required"`
		SampleOutput  string `form:"expected_output" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Edit Question",
			"error": err.Error(),
		})
		return
	}

	var question model.Problem
	userID := sessions.Default(c).Get("user_id")
	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	if err := db.DB.First(&question, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not found",
		})
		return
	}
	if question.OwnerID != userID && !isAdmin {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Question not available",
		})
	}

	updates := map[string]interface{}{
		"title":           form.Title,
		"statement":       form.Statement,
		"time_limit_ms":   form.TimeLimitMs,
		"memory_limit_mb": form.MemoryLimitMs,
		"sample_input":    form.SampleInput,
		"sample_output":   form.SampleOutput,
	}
	if err := db.DB.Model(&question).Updates(updates).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Edit Question",
			"error": "Failed to save question",
		})
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"userID":   userID,
		"isAdmin":  isAdmin,
		"title":    "Edit Question",
		"question": question,
	})
}

// MyQuestionsHandler shows all questions the current user created.
func MyQuestionsHandler(c *gin.Context) {
	userID := sessions.Default(c).Get("user_id")
	var questions []model.Problem
	db.DB.
		Where("owner_id = ?", userID).
		Order("created_at DESC").
		Find(&questions)
	fmt.Println("%s", userID)

	isAdmin := sessions.Default(c).Get("is_admin").(bool)
	c.HTML(http.StatusOK, "my-questions.html", gin.H{
		"userID":    userID,
		"isAdmin":   isAdmin,
		"title":     "My Questions",
		"questions": questions,
	})
}
