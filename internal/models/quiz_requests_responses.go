package models

import (
	"time"
)

// Request Models
type StartQuizRequest struct {
	TopicID        uint            `json:"topic_id" validate:"required"`
	Difficulty     DifficultyLevel `json:"difficulty" validate:"required"`
	TotalQuestions int             `json:"total_questions" validate:"required,min=1,max=50"`
}

type SubmitAnswerRequest struct {
	QuestionID       uint  `json:"question_id" validate:"required"`
	SelectedOptionID *uint `json:"selected_option_id"`
	TimeSpentSeconds int   `json:"time_spent_seconds"`
}

// Response Models
type QuizSessionResponse struct {
	ID             uint            `json:"id"`
	Topic          Topic           `json:"topic"`
	Difficulty     DifficultyLevel `json:"difficulty"`
	TotalQuestions int             `json:"total_questions"`
	Status         SessionStatus   `json:"status"`
	StartedAt      time.Time       `json:"started_at"`
	Progress       QuizProgress    `json:"progress"`
}

type QuizProgress struct {
	Answered  int `json:"answered"`
	Remaining int `json:"remaining"`
}

type QuestionResponse struct {
	ID           uint             `json:"id"`
	QuestionText string           `json:"question_text"`
	Difficulty   DifficultyLevel  `json:"difficulty"`
	Options      []OptionResponse `json:"options"`
	Progress     ProgressInfo     `json:"progress"`
}

type OptionResponse struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
}

type ProgressInfo struct {
	CurrentQuestion int `json:"current_question"`
	TotalQuestions  int `json:"total_questions"`
}

type QuestionOptionResponse struct {
	ID          uint   `json:"id"`
	OptionText  string `json:"option_text"`
	OptionOrder int    `json:"option_order"`
}

type QuizResultsResponse struct {
	SessionID       uint             `json:"session_id"`
	TopicName       string           `json:"topic_name"`
	Difficulty      DifficultyLevel  `json:"difficulty"`
	TotalQuestions  int              `json:"total_questions"`
	CorrectAnswers  int              `json:"correct_answers"`
	TotalPoints     int              `json:"total_points"`
	FinalScore      int              `json:"final_score"`
	TimeSpentTotal  int              `json:"time_spent_total"`
	CompletedAt     time.Time        `json:"completed_at"`
	QuestionResults []QuestionResult `json:"question_results"`
}

type QuestionResult struct {
	QuestionID       uint            `json:"question_id"`
	QuestionText     string          `json:"question_text"`
	Difficulty       DifficultyLevel `json:"difficulty"`
	SelectedOptionID *uint           `json:"selected_option_id"`
	IsCorrect        bool            `json:"is_correct"`
	TimeSpentSeconds int             `json:"time_spent_seconds"`
	Explanation      string          `json:"explanation"`
	Options          []OptionResult  `json:"options"`
}

type OptionResult struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	IsSelected bool   `json:"is_selected"`
}

type QuizAnswerResult struct {
	Question       Question       `json:"question"`
	SelectedOption QuestionOption `json:"selected_option"`
	CorrectOption  QuestionOption `json:"correct_option"`
	IsCorrect      bool           `json:"is_correct"`
	Explanation    string         `json:"explanation"`
}

type QuizSessionDetailsResponse struct {
	Session       QuizSession    `json:"session"`
	AnswerDetails []AnswerDetail `json:"answer_details"`
}

type AnswerDetail struct {
	QuestionID         uint      `json:"question_id"`
	QuestionText       string    `json:"question_text"`
	SelectedOptionID   *uint     `json:"selected_option_id"`
	SelectedOptionText string    `json:"selected_option_text,omitempty"`
	IsCorrect          bool      `json:"is_correct"`
	TimeSpentSeconds   int       `json:"time_spent_seconds"`
	AnsweredAt         time.Time `json:"answered_at"`
}

type UserQuizStatsResponse struct {
	TotalSessions          int          `json:"total_sessions"`
	CompletedSessions      int          `json:"completed_sessions"`
	AverageScore           float64      `json:"average_score"`
	TotalQuestionsAnswered int          `json:"total_questions_answered"`
	TotalCorrectAnswers    int          `json:"total_correct_answers"`
	AccuracyPercentage     float64      `json:"accuracy_percentage"`
	TopicStats             []TopicStats `json:"topic_stats"`
}

type TopicStats struct {
	TopicID   uint    `json:"topic_id"`
	TopicName string  `json:"topic_name"`
	Sessions  int64   `json:"sessions"`
	AvgScore  float64 `json:"avg_score"`
}

type SubmitAnswerResponse struct {
	IsCorrect       bool         `json:"is_correct"`
	AnsweredAt      time.Time    `json:"answered_at"`
	Progress        QuizProgress `json:"progress"`
	SessionComplete bool         `json:"session_complete"`
}

type TopicStatsResponse struct {
	TopicID    uint           `json:"topic_id"`
	TopicName  string         `json:"topic_name"`
	Statistics map[string]int `json:"statistics"` // "easy": 10, "intermediate": 15, "hard": 5
}

type QuizHistoryResponse struct {
	Sessions []QuizSessionSummary `json:"sessions"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
}

type QuizSessionSummary struct {
	ID          uint       `json:"id"`
	TopicName   string     `json:"topic_name"`
	Difficulty  string     `json:"difficulty"`
	FinalScore  *int       `json:"final_score"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
