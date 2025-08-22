package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"quizer/datasource"
	"quizer/internal/models"

	"gorm.io/gorm"
)

// QuizService defines the interface for quiz-related operations
type QuizService interface {
	// Topic management
	GetActiveTopics() ([]models.Topic, error)
	GetTopicByID(id uint) (*models.Topic, error)
	GetTopicQuestionStats(topicID uint) (map[string]int, error)

	// Quiz session management
	StartQuizSession(userID uint, req *models.StartQuizRequest) (*models.QuizSession, error)
	GetQuizSession(sessionID, userID uint) (*models.QuizSession, error)
	AbandonQuizSession(sessionID, userID uint) error

	// Question management
	GetNextQuestion(sessionID, userID uint) (*models.QuestionResponse, error)
	SubmitAnswer(sessionID, userID uint, req *models.SubmitAnswerRequest) (*models.QuizAnswer, error)

	// Results management
	CompleteQuizSession(sessionID, userID uint) (*models.QuizResultsResponse, error)
	GetQuizResults(sessionID, userID uint) (*models.QuizResultsResponse, error)

	// User history
	GetUserQuizHistory(userID uint, limit, offset int) ([]models.QuizSession, error)
	GetUserSessionDetails(sessionID, userID uint) (*models.QuizSessionDetailsResponse, error)

	// Statistics
	GetUserQuizStats(userID uint) (*models.UserQuizStatsResponse, error)
}

type quizService struct {
	db *gorm.DB
}

// NewQuizService creates a new quiz service instance
func NewQuizService() QuizService {
	return &quizService{
		db: datasource.GetDB(),
	}
}

// GetActiveTopics retrieves all active topics
func (s *quizService) GetActiveTopics() ([]models.Topic, error) {
	var topics []models.Topic
	err := s.db.Where("is_active = ?", true).Find(&topics).Error
	return topics, err
}

// GetTopicByID retrieves a topic by ID
func (s *quizService) GetTopicByID(id uint) (*models.Topic, error) {
	var topic models.Topic
	err := s.db.Where("id = ? AND is_active = ?", id, true).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// GetTopicQuestionStats returns question count by difficulty for a topic
func (s *quizService) GetTopicQuestionStats(topicID uint) (map[string]int, error) {
	type DifficultyCount struct {
		Difficulty string
		Count      int64
	}

	var results []DifficultyCount
	err := s.db.Model(&models.Question{}).
		Select("difficulty, count(*) as count").
		Where("topic_id = ? AND is_active = ?", topicID, true).
		Group("difficulty").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, result := range results {
		stats[result.Difficulty] = int(result.Count)
	}

	return stats, nil
}

// StartQuizSession creates a new quiz session
func (s *quizService) StartQuizSession(userID uint, req *models.StartQuizRequest) (*models.QuizSession, error) {
	// Check if topic exists and is active
	_, err := s.GetTopicByID(req.TopicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	// Check if there are enough questions for the requested difficulty
	var questionCount int64
	err = s.db.Model(&models.Question{}).
		Where("topic_id = ? AND difficulty = ? AND is_active = ?", req.TopicID, req.Difficulty, true).
		Count(&questionCount).Error
	if err != nil {
		return nil, err
	}

	if int(questionCount) < req.TotalQuestions {
		return nil, fmt.Errorf("not enough questions available. requested: %d, available: %d", req.TotalQuestions, questionCount)
	}

	// Check if user has an active session for this topic
	var activeSession models.QuizSession
	err = s.db.Where("user_id = ? AND topic_id = ? AND status = ?", userID, req.TopicID, models.SessionStatusActive).First(&activeSession).Error
	if err == nil {
		return nil, errors.New("you already have an active session for this topic")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new quiz session
	session := models.QuizSession{
		UserID:         userID,
		TopicID:        req.TopicID,
		Difficulty:     req.Difficulty,
		TotalQuestions: req.TotalQuestions,
		Status:         models.SessionStatusActive,
		StartedAt:      time.Now(),
	}

	err = s.db.Create(&session).Error
	if err != nil {
		return nil, err
	}

	// Load associations
	err = s.db.Preload("Topic").Preload("User").First(&session, session.ID).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetQuizSession retrieves a quiz session by ID and user ID
func (s *quizService) GetQuizSession(sessionID, userID uint) (*models.QuizSession, error) {
	var session models.QuizSession
	err := s.db.Preload("Topic").Preload("User").Preload("Answers").
		Where("id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// AbandonQuizSession marks a quiz session as abandoned
func (s *quizService) AbandonQuizSession(sessionID, userID uint) error {
	result := s.db.Model(&models.QuizSession{}).
		Where("id = ? AND user_id = ? AND status = ?", sessionID, userID, models.SessionStatusActive).
		Update("status", models.SessionStatusAbandoned)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("session not found or already completed/abandoned")
	}

	return nil
}

// GetNextQuestion retrieves the next unanswered question for a quiz session
func (s *quizService) GetNextQuestion(sessionID, userID uint) (*models.QuestionResponse, error) {
	// Verify session belongs to user and is active
	session, err := s.GetQuizSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	if session.Status != models.SessionStatusActive {
		return nil, errors.New("quiz session is not active")
	}

	// Get answered question IDs
	var answeredQuestionIDs []uint
	err = s.db.Model(&models.QuizAnswer{}).
		Where("quiz_session_id = ?", sessionID).
		Pluck("question_id", &answeredQuestionIDs).Error
	if err != nil {
		return nil, err
	}

	// Check if quiz is already completed
	if len(answeredQuestionIDs) >= session.TotalQuestions {
		return nil, errors.New("all questions have been answered")
	}

	// Get questions that user has already answered before (to avoid duplicates)
	var userAnsweredQuestionIDs []uint
	err = s.db.Model(&models.UserQuestionHistory{}).
		Where("user_id = ?", userID).
		Pluck("question_id", &userAnsweredQuestionIDs).Error
	if err != nil {
		return nil, err
	}

	// Build exclusion list (session answers + user history)
	excludeIDs := append(answeredQuestionIDs, userAnsweredQuestionIDs...)

	// Get available questions
	query := s.db.Model(&models.Question{}).
		Where("topic_id = ? AND difficulty = ? AND is_active = ?", session.TopicID, session.Difficulty, true)

	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	var availableQuestions []models.Question
	err = query.Find(&availableQuestions).Error
	if err != nil {
		return nil, err
	}

	if len(availableQuestions) == 0 {
		return nil, errors.New("no more unique questions available")
	}

	// Randomly select a question
	rand.Seed(time.Now().UnixNano())
	selectedQuestion := availableQuestions[rand.Intn(len(availableQuestions))]

	// Load question options
	err = s.db.Preload("Options").First(&selectedQuestion, selectedQuestion.ID).Error
	if err != nil {
		return nil, err
	}

	// Create response without revealing correct answer
	response := &models.QuestionResponse{
		ID:           selectedQuestion.ID,
		QuestionText: selectedQuestion.QuestionText,
		Difficulty:   selectedQuestion.Difficulty,
		Options:      make([]models.OptionResponse, len(selectedQuestion.Options)),
		Progress: models.ProgressInfo{
			CurrentQuestion: len(answeredQuestionIDs) + 1,
			TotalQuestions:  session.TotalQuestions,
		},
	}

	// Shuffle options for better UX
	rand.Shuffle(len(selectedQuestion.Options), func(i, j int) {
		selectedQuestion.Options[i], selectedQuestion.Options[j] = selectedQuestion.Options[j], selectedQuestion.Options[i]
	})

	for i, option := range selectedQuestion.Options {
		response.Options[i] = models.OptionResponse{
			ID:         option.ID,
			OptionText: option.OptionText,
		}
	}

	return response, nil
}

// SubmitAnswer processes a user's answer to a question
func (s *quizService) SubmitAnswer(sessionID, userID uint, req *models.SubmitAnswerRequest) (*models.QuizAnswer, error) {
	// Verify session belongs to user and is active
	session, err := s.GetQuizSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	if session.Status != models.SessionStatusActive {
		return nil, errors.New("quiz session is not active")
	}

	// Check if question has already been answered in this session
	var existingAnswer models.QuizAnswer
	err = s.db.Where("quiz_session_id = ? AND question_id = ?", sessionID, req.QuestionID).First(&existingAnswer).Error
	if err == nil {
		return nil, errors.New("question has already been answered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get the question and verify it belongs to the session topic/difficulty
	var question models.Question
	err = s.db.Preload("Options").Where("id = ? AND topic_id = ? AND difficulty = ? AND is_active = ?",
		req.QuestionID, session.TopicID, session.Difficulty, true).First(&question).Error
	if err != nil {
		return nil, fmt.Errorf("question not found or invalid: %w", err)
	}

	// Find the selected option and check if it's correct
	var selectedOption *models.QuestionOption
	var isCorrect bool

	if req.SelectedOptionID != nil {
		for _, option := range question.Options {
			if option.ID == *req.SelectedOptionID {
				selectedOption = &option
				isCorrect = option.IsCorrect
				break
			}
		}

		if selectedOption == nil {
			return nil, errors.New("selected option not found")
		}
	}

	// Create quiz answer
	answer := models.QuizAnswer{
		QuizSessionID:    sessionID,
		QuestionID:       req.QuestionID,
		SelectedOptionID: req.SelectedOptionID,
		IsCorrect:        isCorrect,
		TimeSpentSeconds: req.TimeSpentSeconds,
		AnsweredAt:       time.Now(),
	}

	err = s.db.Create(&answer).Error
	if err != nil {
		return nil, err
	}

	// Update user question history
	history := models.UserQuestionHistory{
		UserID:            userID,
		QuestionID:        req.QuestionID,
		TimesAnswered:     1,
		TimesCorrect:      0,
		LastAnsweredAt:    time.Now(),
		LastAnswerCorrect: isCorrect,
	}

	if isCorrect {
		history.TimesCorrect = 1
	}

	// Use ON CONFLICT to update existing records
	err = s.db.Create(&history).Error
	if err != nil {
		// If record exists, update it
		err = s.db.Model(&models.UserQuestionHistory{}).
			Where("user_id = ? AND question_id = ?", userID, req.QuestionID).
			Updates(map[string]interface{}{
				"times_answered":      gorm.Expr("times_answered + 1"),
				"times_correct":       gorm.Expr("CASE WHEN ? THEN times_correct + 1 ELSE times_correct END", isCorrect),
				"last_answered_at":    time.Now(),
				"last_answer_correct": isCorrect,
			}).Error
		if err != nil {
			// Log error but don't fail the answer submission
			fmt.Printf("Warning: Failed to update user question history: %v\n", err)
		}
	}

	// Load associations for response
	err = s.db.Preload("Question").Preload("SelectedOption").First(&answer, answer.ID).Error
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

// CompleteQuizSession completes a quiz session and calculates final score
func (s *quizService) CompleteQuizSession(sessionID, userID uint) (*models.QuizResultsResponse, error) {
	// Verify session belongs to user and is active
	session, err := s.GetQuizSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	if session.Status != models.SessionStatusActive {
		return nil, errors.New("quiz session is not active")
	}

	// Get all answers for this session
	var answers []models.QuizAnswer
	err = s.db.Preload("Question").Preload("Question.Options").Preload("SelectedOption").
		Where("quiz_session_id = ?", sessionID).
		Find(&answers).Error
	if err != nil {
		return nil, err
	}

	// Calculate score
	correctAnswers := 0
	totalAnswers := len(answers)
	totalPoints := 0

	for _, answer := range answers {
		if answer.IsCorrect {
			correctAnswers++
		}
	}

	// Calculate percentage score
	var finalScore int
	if totalAnswers > 0 {
		finalScore = (correctAnswers * 100) / totalAnswers
	}

	// Update session status
	now := time.Now()
	err = s.db.Model(&session).Updates(map[string]interface{}{
		"status":       models.SessionStatusCompleted,
		"completed_at": &now,
		"final_score":  &finalScore,
	}).Error
	if err != nil {
		return nil, err
	}

	// Prepare detailed results
	questionResults := make([]models.QuestionResult, len(answers))
	for i, answer := range answers {
		questionResults[i] = models.QuestionResult{
			QuestionID:       answer.QuestionID,
			QuestionText:     answer.Question.QuestionText,
			Difficulty:       answer.Question.Difficulty,
			SelectedOptionID: answer.SelectedOptionID,
			IsCorrect:        answer.IsCorrect,
			TimeSpentSeconds: answer.TimeSpentSeconds,
			Explanation:      answer.Question.Explanation,
			Options:          make([]models.OptionResult, len(answer.Question.Options)),
		}

		// Add all options with correct answer indication
		for j, option := range answer.Question.Options {
			questionResults[i].Options[j] = models.OptionResult{
				ID:         option.ID,
				OptionText: option.OptionText,
				IsCorrect:  option.IsCorrect,
				IsSelected: answer.SelectedOptionID != nil && option.ID == *answer.SelectedOptionID,
			}
		}
	}

	// Calculate time spent
	timeSpent := int(now.Sub(session.StartedAt).Seconds())

	response := &models.QuizResultsResponse{
		SessionID:       sessionID,
		TopicName:       session.Topic.Name,
		Difficulty:      session.Difficulty,
		TotalQuestions:  session.TotalQuestions,
		CorrectAnswers:  correctAnswers,
		TotalPoints:     totalPoints,
		FinalScore:      finalScore,
		TimeSpentTotal:  timeSpent,
		CompletedAt:     now,
		QuestionResults: questionResults,
	}

	return response, nil
}

// GetQuizResults retrieves results for a completed quiz session
func (s *quizService) GetQuizResults(sessionID, userID uint) (*models.QuizResultsResponse, error) {
	// Verify session belongs to user
	session, err := s.GetQuizSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	if session.Status != models.SessionStatusCompleted {
		return nil, errors.New("quiz session is not completed yet")
	}

	// Get all answers for this session
	var answers []models.QuizAnswer
	err = s.db.Preload("Question").Preload("Question.Options").Preload("SelectedOption").
		Where("quiz_session_id = ?", sessionID).
		Find(&answers).Error
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	correctAnswers := 0
	totalPoints := 0

	for _, answer := range answers {
		if answer.IsCorrect {
			correctAnswers++
		}
	}

	// Prepare detailed results
	questionResults := make([]models.QuestionResult, len(answers))
	for i, answer := range answers {
		questionResults[i] = models.QuestionResult{
			QuestionID:       answer.QuestionID,
			QuestionText:     answer.Question.QuestionText,
			Difficulty:       answer.Question.Difficulty,
			SelectedOptionID: answer.SelectedOptionID,
			IsCorrect:        answer.IsCorrect,
			TimeSpentSeconds: answer.TimeSpentSeconds,
			Explanation:      answer.Question.Explanation,
			Options:          make([]models.OptionResult, len(answer.Question.Options)),
		}

		// Add all options with correct answer indication
		for j, option := range answer.Question.Options {
			questionResults[i].Options[j] = models.OptionResult{
				ID:         option.ID,
				OptionText: option.OptionText,
				IsCorrect:  option.IsCorrect,
				IsSelected: answer.SelectedOptionID != nil && option.ID == *answer.SelectedOptionID,
			}
		}
	}

	// Calculate time spent
	var timeSpent int
	if session.CompletedAt != nil {
		timeSpent = int(session.CompletedAt.Sub(session.StartedAt).Seconds())
	}

	response := &models.QuizResultsResponse{
		SessionID:       sessionID,
		TopicName:       session.Topic.Name,
		Difficulty:      session.Difficulty,
		TotalQuestions:  session.TotalQuestions,
		CorrectAnswers:  correctAnswers,
		TotalPoints:     totalPoints,
		FinalScore:      *session.FinalScore,
		TimeSpentTotal:  timeSpent,
		CompletedAt:     *session.CompletedAt,
		QuestionResults: questionResults,
	}

	return response, nil
}

// GetUserQuizHistory retrieves user's quiz session history
func (s *quizService) GetUserQuizHistory(userID uint, limit, offset int) ([]models.QuizSession, error) {
	var sessions []models.QuizSession

	query := s.db.Preload("Topic").Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&sessions).Error
	return sessions, err
}

// GetUserSessionDetails retrieves detailed information about a specific session
func (s *quizService) GetUserSessionDetails(sessionID, userID uint) (*models.QuizSessionDetailsResponse, error) {
	// Get session
	session, err := s.GetQuizSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	// Get session answers with questions
	var answers []models.QuizAnswer
	err = s.db.Preload("Question").Preload("SelectedOption").
		Where("quiz_session_id = ?", sessionID).
		Order("answered_at ASC").
		Find(&answers).Error
	if err != nil {
		return nil, err
	}

	// Prepare answer details
	answerDetails := make([]models.AnswerDetail, len(answers))
	for i, answer := range answers {
		answerDetails[i] = models.AnswerDetail{
			QuestionID:       answer.QuestionID,
			QuestionText:     answer.Question.QuestionText,
			SelectedOptionID: answer.SelectedOptionID,
			IsCorrect:        answer.IsCorrect,
			TimeSpentSeconds: answer.TimeSpentSeconds,
			AnsweredAt:       answer.AnsweredAt,
		}

		if answer.SelectedOption != nil {
			answerDetails[i].SelectedOptionText = answer.SelectedOption.OptionText
		}
	}

	response := &models.QuizSessionDetailsResponse{
		Session:       *session,
		AnswerDetails: answerDetails,
	}

	return response, nil
}

// GetUserQuizStats calculates and returns user's quiz statistics
func (s *quizService) GetUserQuizStats(userID uint) (*models.UserQuizStatsResponse, error) {
	// Total sessions
	var totalSessions int64
	err := s.db.Model(&models.QuizSession{}).Where("user_id = ?", userID).Count(&totalSessions).Error
	if err != nil {
		return nil, err
	}

	// Completed sessions
	var completedSessions int64
	err = s.db.Model(&models.QuizSession{}).
		Where("user_id = ? AND status = ?", userID, models.SessionStatusCompleted).
		Count(&completedSessions).Error
	if err != nil {
		return nil, err
	}

	// Average score
	var avgScore float64
	err = s.db.Model(&models.QuizSession{}).
		Where("user_id = ? AND status = ? AND final_score IS NOT NULL", userID, models.SessionStatusCompleted).
		Select("AVG(final_score)").Scan(&avgScore).Error
	if err != nil {
		return nil, err
	}

	// Total questions answered
	var totalQuestionsAnswered int64
	err = s.db.Table("quiz_answers").
		Joins("JOIN quiz_sessions ON quiz_answers.quiz_session_id = quiz_sessions.id").
		Where("quiz_sessions.user_id = ?", userID).
		Count(&totalQuestionsAnswered).Error
	if err != nil {
		return nil, err
	}

	// Total correct answers
	var totalCorrectAnswers int64
	err = s.db.Table("quiz_answers").
		Joins("JOIN quiz_sessions ON quiz_answers.quiz_session_id = quiz_sessions.id").
		Where("quiz_sessions.user_id = ? AND quiz_answers.is_correct = ?", userID, true).
		Count(&totalCorrectAnswers).Error
	if err != nil {
		return nil, err
	}

	// Accuracy percentage
	var accuracyPercentage float64
	if totalQuestionsAnswered > 0 {
		accuracyPercentage = (float64(totalCorrectAnswers) / float64(totalQuestionsAnswered)) * 100
	}

	// Stats by topic
	var topicStats []models.TopicStats
	err = s.db.Table("quiz_sessions").
		Select("quiz_sessions.topic_id, topics.name as topic_name, COUNT(*) as sessions, AVG(quiz_sessions.final_score) as avg_score").
		Joins("JOIN topics ON quiz_sessions.topic_id = topics.id").
		Where("quiz_sessions.user_id = ? AND quiz_sessions.status = ?", userID, models.SessionStatusCompleted).
		Group("quiz_sessions.topic_id, topics.name").
		Scan(&topicStats).Error
	if err != nil {
		return nil, err
	}

	response := &models.UserQuizStatsResponse{
		TotalSessions:          int(totalSessions),
		CompletedSessions:      int(completedSessions),
		AverageScore:           avgScore,
		TotalQuestionsAnswered: int(totalQuestionsAnswered),
		TotalCorrectAnswers:    int(totalCorrectAnswers),
		AccuracyPercentage:     accuracyPercentage,
		TopicStats:             topicStats,
	}

	return response, nil
}
