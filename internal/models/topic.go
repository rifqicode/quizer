package models

type Topic struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Associations
	Questions []Question `gorm:"foreignKey:TopicID" json:"questions,omitempty"`
}

func (Topic) TableName() string {
	return "topics"
}
