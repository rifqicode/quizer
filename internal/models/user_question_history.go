package models

import (
	"time"
)

type UserQuestionHistory struct {
	BaseModel
	UserID            uint      `gorm:"not null;index:idx_user_question_history_user_question,unique" json:"user_id"`
	QuestionID        uint      `gorm:"not null;index:idx_user_question_history_user_question,unique" json:"question_id"`
	TimesAnswered     int       `gorm:"default:0" json:"times_answered"`
	TimesCorrect      int       `gorm:"default:0" json:"times_correct"`
	LastAnsweredAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_answered_at"`
	LastAnswerCorrect bool      `gorm:"index" json:"last_answer_correct"`

	// Associations
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Question Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (UserQuestionHistory) TableName() string {
	return "user_question_history"
}
