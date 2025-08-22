package models

type DifficultyLevel string

const (
	DifficultyEasy         DifficultyLevel = "easy"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyHard         DifficultyLevel = "hard"
)

type Question struct {
	BaseModel
	TopicID      uint            `gorm:"not null;index" json:"topic_id"`
	QuestionText string          `gorm:"type:text;not null" json:"question_text"`
	Difficulty   DifficultyLevel `gorm:"type:varchar(20);not null;check:difficulty IN ('easy','intermediate','hard');index" json:"difficulty"`
	Explanation  string          `gorm:"type:text;not null" json:"explanation"`
	IsActive     bool            `gorm:"default:true;index" json:"is_active"`

	// Associations
	Topic   Topic            `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Options []QuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

func (Question) TableName() string {
	return "questions"
}
