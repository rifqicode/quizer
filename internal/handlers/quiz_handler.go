package handlers

import (
	"net/http"
	"strconv"

	"quizer/internal/models"
	"quizer/internal/services"
	"quizer/pkg/utils"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService services.QuizService
}

// NewQuizHandler creates a new quiz handler instance
func NewQuizHandler(quizService services.QuizService) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

// GetActiveTopics godoc
// @Summary Get all active topics
// @Description Retrieve all active topics available for quizzes
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.Topic}
// @Failure 500 {object} utils.Response
// @Router /api/v1/quiz/topics [get]
func (h *QuizHandler) GetActiveTopics(c *gin.Context) {
	topics, err := h.quizService.GetActiveTopics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve topics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Topics retrieved successfully", topics)
}

// GetTopicByID godoc
// @Summary Get topic by ID
// @Description Retrieve a specific topic by its ID
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Success 200 {object} utils.Response{data=models.Topic}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/quiz/topics/{id} [get]
func (h *QuizHandler) GetTopicByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid topic ID")
		return
	}

	topic, err := h.quizService.GetTopicByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Topic not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Topic retrieved successfully", topic)
}

// GetTopicQuestionStats godoc
// @Summary Get topic question statistics
// @Description Retrieve question count by difficulty for a topic
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Success 200 {object} utils.Response{data=map[string]int}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/quiz/topics/{id}/stats [get]
func (h *QuizHandler) GetTopicQuestionStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid topic ID")
		return
	}

	stats, err := h.quizService.GetTopicQuestionStats(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve topic statistics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Topic statistics retrieved successfully", stats)
}

// StartQuizSession godoc
// @Summary Start a new quiz session
// @Description Create a new quiz session for the authenticated user
// @Tags quiz-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.StartQuizRequest true "Start quiz request"
// @Success 201 {object} utils.Response{data=models.QuizSession}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/quiz/sessions [post]
func (h *QuizHandler) StartQuizSession(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.StartQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.TopicID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Topic ID is required")
		return
	}

	if req.TotalQuestions <= 0 || req.TotalQuestions > 50 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Total questions must be between 1 and 50")
		return
	}

	validDifficulties := map[models.DifficultyLevel]bool{
		models.DifficultyEasy:         true,
		models.DifficultyIntermediate: true,
		models.DifficultyHard:         true,
	}

	if !validDifficulties[req.Difficulty] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid difficulty level")
		return
	}

	session, err := h.quizService.StartQuizSession(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to start quiz session")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Quiz session started successfully", session)
}

// GetQuizSession godoc
// @Summary Get quiz session details
// @Description Retrieve details of a specific quiz session
// @Tags quiz-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} utils.Response{data=models.QuizSession}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id} [get]
func (h *QuizHandler) GetQuizSession(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := h.quizService.GetQuizSession(uint(sessionID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Session not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Session retrieved successfully", session)
}

// AbandonQuizSession godoc
// @Summary Abandon a quiz session
// @Description Mark a quiz session as abandoned
// @Tags quiz-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id}/abandon [put]
func (h *QuizHandler) AbandonQuizSession(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	err = h.quizService.AbandonQuizSession(uint(sessionID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to abandon session")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Session abandoned successfully", nil)
}

// GetNextQuestion godoc
// @Summary Get next question in quiz session
// @Description Retrieve the next unanswered question for a quiz session
// @Tags questions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} utils.Response{data=models.QuestionResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id}/next-question [get]
func (h *QuizHandler) GetNextQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	question, err := h.quizService.GetNextQuestion(uint(sessionID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get next question")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Next question retrieved successfully", question)
}

// SubmitAnswer godoc
// @Summary Submit answer for a question
// @Description Submit user's answer for a specific question in a quiz session
// @Tags questions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Param request body models.SubmitAnswerRequest true "Submit answer request"
// @Success 200 {object} utils.Response{data=models.QuizAnswer}
// @Failure 400 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id}/answer [post]
func (h *QuizHandler) SubmitAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	var req models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	answer, err := h.quizService.SubmitAnswer(uint(sessionID), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to submit answer")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Answer submitted successfully", answer)
}

// CompleteQuizSession godoc
// @Summary Complete quiz session
// @Description Complete a quiz session and get final results
// @Tags quiz-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} utils.Response{data=models.QuizResultsResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id}/complete [post]
func (h *QuizHandler) CompleteQuizSession(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	results, err := h.quizService.CompleteQuizSession(uint(sessionID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to complete quiz session")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz session completed successfully", results)
}

// GetQuizResults godoc
// @Summary Get quiz results
// @Description Retrieve results for a completed quiz session
// @Tags quiz-sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} utils.Response{data=models.QuizResultsResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/quiz/sessions/{id}/results [get]
func (h *QuizHandler) GetQuizResults(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	results, err := h.quizService.GetQuizResults(uint(sessionID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get quiz results")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz results retrieved successfully", results)
}

// GetUserQuizHistory godoc
// @Summary Get user quiz history
// @Description Retrieve user's quiz session history
// @Tags quiz-history
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} utils.Response{data=[]models.QuizSession}
// @Failure 400 {object} utils.Response
// @Router /api/v1/quiz/my-sessions [get]
func (h *QuizHandler) GetUserQuizHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	sessions, err := h.quizService.GetUserQuizHistory(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve quiz history")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz history retrieved successfully", sessions)
}

// GetUserQuizStats godoc
// @Summary Get user quiz statistics
// @Description Retrieve user's quiz statistics and performance metrics
// @Tags quiz-stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserQuizStatsResponse}
// @Failure 500 {object} utils.Response
// @Router /api/v1/quiz/my-stats [get]
func (h *QuizHandler) GetUserQuizStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.quizService.GetUserQuizStats(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve quiz statistics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz statistics retrieved successfully", stats)
}
