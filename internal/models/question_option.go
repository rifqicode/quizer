package models

type QuestionOption struct {
	BaseModel
	QuestionID uint   `gorm:"not null;index" json:"question_id"`
	OptionText string `gorm:"type:text;not null" json:"option_text"`
	IsCorrect  bool   `gorm:"default:false" json:"is_correct"`

	// Associations
	Question Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (QuestionOption) TableName() string {
	return "question_options"
}
