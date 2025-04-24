package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// MySubmissionsHandler displays the user's submission history
func MySubmissionsHandler(c *gin.Context) {
	// For Phase 1, generate mock submissions
	submissions := getMockSubmissions()

	c.HTML(http.StatusOK, "sList.html", gin.H{
		"title":       "My Submissions",
		"submissions": submissions,
	})
}

// SubmissionDetailHandler displays the details of a submission
func SubmissionDetailHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Not Found",
			"error": "Submission not found",
		})
		return
	}

	// For Phase 1, generate a mock submission
	submission := getMockSubmission(uint(id))

	c.HTML(http.StatusOK, "submission/sDetail.html", gin.H{
		"title":      "Submission Details",
		"submission": submission,
	})
}

// SubmitHandler handles code submission for a question
func SubmitHandler(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("question_id"))
	if err != nil || questionID < 1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Not Found",
			"error": "Question not found",
		})
		return
	}

	// For Phase 1, just redirect back to the question
	// We'll implement actual submission in later phases
	c.Redirect(http.StatusSeeOther, "/questions/"+strconv.Itoa(questionID))
}

// Mock data helpers

func getMockSubmissions() []model.Submission {
	statuses := []string{
		model.StatusOK,
		model.StatusCompileError,
		model.StatusWrongAnswer,
		model.StatusMemoryLimitExceeded,
		model.StatusTimeLimitExceeded,
		model.StatusRuntimeError,
		model.StatusPending,
	}

	submissions := make([]model.Submission, 10)
	for i := 0; i < 10; i++ {
		submissions[i] = model.Submission{
			ID:        uint(i + 1),
			ProblemID: uint((i % 5) + 1),
			UserID:    1,
			Code:      "package main\n\nfunc main() {\n\t// Some code\n}",
			Language:  "go",
			Status:    statuses[i%len(statuses)],
			CreatedAt: time.Now().AddDate(0, 0, -i),
			UpdatedAt: time.Now().AddDate(0, 0, -i).Add(time.Second * 30),
		}
	}
	return submissions
}

// internal/handler/submission.go
func getMockSubmission(id uint) model.Submission {
	statuses := []string{
		model.StatusOK,
		model.StatusCompileError,
		model.StatusWrongAnswer,
		model.StatusMemoryLimitExceeded,
		model.StatusTimeLimitExceeded,
		model.StatusRuntimeError,
		model.StatusPending,
	}

	return model.Submission{
		ID:        id,
		ProblemID: uint((int(id) % 5) + 1),
		UserID:    1,
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scan(&a, &b)\n\tfmt.Println(a + b)\n}",
		Language:  "go",
		Status:    statuses[int(id)%len(statuses)],
		CreatedAt: time.Now().AddDate(0, 0, -int(id)),
		UpdatedAt: time.Now().AddDate(0, 0, -int(id)).Add(time.Second * 30),
	}
}
