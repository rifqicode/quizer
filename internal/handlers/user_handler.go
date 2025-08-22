package handlers

import (
	"log"
	"net/http"
	"quizer/internal/models"
	"quizer/internal/services"
	"quizer/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in CreateUser: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	var user models.User

	if err := utils.BindJSON(c, &user); err != nil {
		return // BindJSON already sends error response
	}

	if err := h.userService.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in GetUser: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		log.Printf("Error getting user by ID %d: %v", id, err)
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users from the database
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in GetAllUsers: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	users, err := h.userService.GetAllUsers()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in UpdateUser: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var user models.User
	if err := utils.BindJSON(c, &user); err != nil {
		return // BindJSON already sends error response
	}

	if err := h.userService.UpdateUser(uint(id), &user); err != nil {
		log.Printf("Error updating user ID %d: %v", id, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in DeleteUser: %v", r)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		log.Printf("Error deleting user ID %d: %v", id, err)
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}
