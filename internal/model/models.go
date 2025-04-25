package model

import "time"

const (
	StatusOK                  = "ok"
	StatusCompileError        = "compile_error"
	StatusWrongAnswer         = "wrong_answer"
	StatusMemoryLimitExceeded = "memory_limit"
	StatusTimeLimitExceeded   = "time_limit"
	StatusRuntimeError        = "runtime_error"
	StatusPending             = "pending"
)

// User represents a platform user.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:50;unique;not null"`
	PasswordHash string `gorm:"not null;default:''"`
	Email        string `gorm:"size:100;unique;not null"`
	IsAdmin      bool   `gorm:"default:false;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Problem is a coding challenge.
type Problem struct {
	ID            uint       `gorm:"primaryKey"`
	OwnerID       uint       `gorm:"not null;index"`
	Title         string     `gorm:"size:200;not null"`
	Statement     string     `gorm:"type:text;not null"`
	TimeLimitMs   int        `gorm:"not null"`
	MemoryLimitMb int        `gorm:"not null"`
	Status        string     `gorm:"size:20;default:'draft';index:idx_status_publish,priority:1"`
	PublishDate   *time.Time `gorm:"index:idx_status_publish,priority:2"`
	SampleInput   string     `gorm:"type:text;not null"`
	SampleOutput  string     `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Submission is a code submission for a problem.
type Submission struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index:idx_user_created,priority:1"`
	ProblemID    uint      `gorm:"not null;index:idx_problem_status,priority:1"`
	Code         string    `gorm:"type:text;not null"`
	Language     string    `gorm:"size:20;not null"`
	Status       string    `gorm:"size:20;default:'pending';index:idx_problem_status,priority:2"`
	CompileError string    `gorm:"type:text"`
	RuntimeError string    `gorm:"type:text"`
	Output       string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"index:idx_user_created,priority:2"`
	UpdatedAt    time.Time
}

// Session tracks a user login session.
type Session struct {
	Token     string    `gorm:"primaryKey;size:128"`
	UserID    uint      `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

// TestCase holds an individual test case for a problem.
type TestCase struct {
	ID             uint   `gorm:"primaryKey"`
	ProblemID      uint   `gorm:"not null;index"`
	InputData      string `gorm:"type:text;not null"`
	ExpectedOutput string `gorm:"type:text;not null"`
	CreatedAt      time.Time
}
