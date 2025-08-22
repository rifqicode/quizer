package models

import (
	"time"
)

type QuizAnswer struct {
	BaseModel
	QuizSessionID    uint      `gorm:"not null;index" json:"quiz_session_id"`
	QuestionID       uint      `gorm:"not null;index" json:"question_id"`
	SelectedOptionID *uint     `gorm:"index" json:"selected_option_id"`
	IsCorrect        bool      `gorm:"index" json:"is_correct"`
	TimeSpentSeconds int       `gorm:"default:0" json:"time_spent_seconds"`
	AnsweredAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"answered_at"`

	// Associations
	QuizSession    QuizSession     `gorm:"foreignKey:QuizSessionID" json:"quiz_session,omitempty"`
	Question       Question        `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	SelectedOption *QuestionOption `gorm:"foreignKey:SelectedOptionID" json:"selected_option,omitempty"`
}

func (QuizAnswer) TableName() string {
	return "quiz_answers"
}
