package models

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusAbandoned SessionStatus = "abandoned"
)

type QuizSession struct {
	BaseModel
	UserID         uint            `gorm:"not null;index" json:"user_id"`
	TopicID        uint            `gorm:"not null;index" json:"topic_id"`
	Difficulty     DifficultyLevel `gorm:"type:varchar(20);not null;check:difficulty IN ('easy','intermediate','hard');index" json:"difficulty"`
	TotalQuestions int             `gorm:"not null" json:"total_questions"`
	Status         SessionStatus   `gorm:"type:varchar(20);default:'active';check:status IN ('active','completed','abandoned');index" json:"status"`
	StartedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	FinalScore     *int            `json:"final_score"`

	// Associations
	User    User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Topic   Topic        `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Answers []QuizAnswer `gorm:"foreignKey:QuizSessionID" json:"answers,omitempty"`
}

func (QuizSession) TableName() string {
	return "quiz_sessions"
}
