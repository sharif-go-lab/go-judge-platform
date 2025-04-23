package model

import "time"

// User represents a user in the system
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password hash, never returned in JSON
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// Question represents a coding problem
type Question struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Statement      string    `json:"statement"`
	TimeLimit      int       `json:"time_limit"`   // in milliseconds
	MemoryLimit    int       `json:"memory_limit"` // in megabytes
	InputTest      string    `json:"input_test"`
	ExpectedOutput string    `json:"expected_output"`
	OwnerID        uint      `json:"owner_id"`
	IsPublished    bool      `json:"is_published"`
	PublishedAt    time.Time `json:"published_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Submission represents a code submission for a question
type Submission struct {
	ID          uint      `json:"id"`
	QuestionID  uint      `json:"question_id"`
	UserID      uint      `json:"user_id"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Status      string    `json:"status"` // Pending, OK, Compile Error, etc.
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

// Status constants for submissions
const (
	StatusPending             = "Pending"
	StatusInProgress          = "In Progress"
	StatusOK                  = "OK"
	StatusCompileError        = "Compile Error"
	StatusWrongAnswer         = "Wrong Answer"
	StatusMemoryLimitExceeded = "Memory Limit Exceeded"
	StatusTimeLimitExceeded   = "Time Limit Exceeded"
	StatusRuntimeError        = "Runtime Error"
)
