package models

type Question struct {
	BaseModel
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	IsActive bool   `json:"is_active" gorm:"default:true"`
}
