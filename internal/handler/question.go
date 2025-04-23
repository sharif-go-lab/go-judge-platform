package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// QuestionListHandler displays the list of questions with pagination
func QuestionListHandler(c *gin.Context) {
	// Get page from query params, default to 1
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// For Phase 1, generate some mock questions
	questions := getMockQuestions()

	// Calculate pagination
	perPage := 10
	totalPages := (len(questions) + perPage - 1) / perPage

	startIdx := (page - 1) * perPage
	endIdx := startIdx + perPage
	if endIdx > len(questions) {
		endIdx = len(questions)
	}

	c.HTML(http.StatusOK, "qList.html", gin.H{
		"title":      "Questions",
		"questions":  questions[startIdx:endIdx],
		"page":       page,
		"totalPages": totalPages,
		"prevPage":   page - 1,
		"nextPage":   page + 1,
	})
}

// QuestionDetailHandler displays a single question
func QuestionDetailHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Not Found",
			"error": "Question not found",
		})
		return
	}

	// For Phase 1, generate a mock question
	question := getMockQuestion(uint(id))

	c.HTML(http.StatusOK, "qDetail.html", gin.H{
		"title":    question.Title,
		"question": question,
	})
}

// CreateQuestionPageHandler displays the question creation form
func CreateQuestionPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{
		"title": "Create Question",
	})
}

// CreateQuestionHandler handles question creation form submission
func CreateQuestionHandler(c *gin.Context) {
	// For Phase 1, just redirect to the questions page
	// We'll implement actual question creation in later phases
	c.Redirect(http.StatusSeeOther, "/questions")
}

// MyQuestionsHandler displays questions created by the current user
func MyQuestionsHandler(c *gin.Context) {
	// For Phase 1, generate some mock questions
	questions := getMockQuestions()

	c.HTML(http.StatusOK, "my-questions.html", gin.H{
		"title":     "My Questions",
		"questions": questions,
	})
}

// Mock data helpers

func getMockQuestions() []model.Question {
	questions := make([]model.Question, 25)
	for i := 0; i < 25; i++ {
		questions[i] = model.Question{
			ID:          uint(i + 1),
			Title:       "Question " + strconv.Itoa(i+1),
			Statement:   "This is a sample question statement for question " + strconv.Itoa(i+1),
			TimeLimit:   1000,
			MemoryLimit: 256,
			OwnerID:     1,
			IsPublished: true,
			PublishedAt: time.Now().AddDate(0, 0, -i),
			CreatedAt:   time.Now().AddDate(0, 0, -i-1),
		}
	}
	return questions
}

func getMockQuestion(id uint) model.Question {
	return model.Question{
		ID:             id,
		Title:          "Question " + strconv.Itoa(int(id)),
		Statement:      "# Question " + strconv.Itoa(int(id)) + "\n\nWrite a program that adds two numbers.\n\n## Input\nTwo integers a and b (1 ≤ a, b ≤ 10^9).\n\n## Output\nThe sum of a and b.",
		TimeLimit:      1000,
		MemoryLimit:    256,
		InputTest:      "1 2",
		ExpectedOutput: "3",
		OwnerID:        1,
		IsPublished:    true,
		PublishedAt:    time.Now().AddDate(0, 0, -int(id)),
		CreatedAt:      time.Now().AddDate(0, 0, -int(id)-1),
	}
}
