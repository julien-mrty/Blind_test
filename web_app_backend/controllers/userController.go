package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julien-mrty/Web_app_jump_higher/web_app_backend/models"
	"github.com/julien-mrty/Web_app_jump_higher/web_app_backend/services"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Add a new user to the database
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var user models.User
	var err error

	// Bind JSON data to the user struct
	if err = c.ShouldBindJSON(&user); err != nil {
		services.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate user data
	if err = services.ValidateStruct(&user); err != nil {
		services.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the user's password
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPassword

	// Generate RoboHash Avatar URL
	user.AvatarURL = fmt.Sprintf("https://robohash.org/%s?set=set1&size=150x150", user.Username)

	// Attempt to create the user
	err = services.CreateUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			// Return a 409 Conflict if the username is taken
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		// Handle other errors
		services.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the created user
	c.JSON(http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve the list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve a specific user by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := services.GetUserByID(id)
	if err != nil {
		services.HandleError(c, http.StatusNotFound, "User not found")
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by their ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteUserByID(id)
	if err != nil {
		services.HandleError(c, http.StatusNotFound, "User not found")
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Description Update a user's information by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [patch]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		services.HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	updatedUser, err := services.UpdateUserByID(id, &user)
	if err != nil {
		services.HandleError(c, http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// Signup godoc
// @Summary User signup
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Signup data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/signup [post]
func Signup(c *gin.Context) {
	var user models.User
	var err error

	// Bind incoming JSON to user struct
	if err = c.ShouldBindJSON(&user); err != nil {
		services.HandleError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Validate the user input
	if err = services.ValidateStruct(&user); err != nil {
		services.HandleError(c, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	// Check if username already exists
	existingUser, err := services.GetUserByUserName(user.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	} else if err != nil && err.Error() != "user not found" {
		// Handle unexpected errors from the database
		services.HandleError(c, http.StatusInternalServerError, "An unexpected error occurred: "+err.Error())
		return
	}

	// Hash the user's password
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPassword

	// Generate RoboHash Avatar URL
	user.AvatarURL = fmt.Sprintf("https://robohash.org/%s?set=set1&size=150x150", user.Username)

	// Attempt to create the user in the database
	err = services.CreateUser(&user)
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	// Generate a JWT token for the newly created user
	token, err := services.GenerateJWT(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username, // Inclure le username dans le payload
		"avatar":   fmt.Sprintf("https://robohash.org/%s?set=set1&size=150x150", user.Username),
	})
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Respond with success and include the token
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var request models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		services.HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := services.GetUserByUserName(request.Username)
	if err != nil {
		services.HandleError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	isValid := services.CheckPasswordHash(request.Password, user.Password)
	if !isValid {
		services.HandleError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := services.GenerateJWT(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username, // Inclure le username dans le payload
		"avatar":   fmt.Sprintf("https://robohash.org/%s?set=set1&size=150x150", user.Username),
	})
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Include token in the response
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "User logged successfully"})
}
