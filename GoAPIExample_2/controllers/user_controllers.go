package controllers

import (
	"net/http"

	"Module/API/models"

	"github.com/gin-gonic/gin"
)

// Create user input model
type UserInput struct {
	ID       uint   `json:"userid"`
	UserName string `json:"username" binding:"required"`
}

/* GET ALL */
// GET '/user'
// Get all user
func FindAllUser(c *gin.Context) {
	var users []models.User
	var userResponses []UserInput

	models.DB.Find(&users)

	for _, user := range users {
		userResponses = append(userResponses, UserInput{
			ID:       user.ID,
			UserName: user.UserName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": userResponses})
}

/* GET */
// GET '/user/todolist'
// Get all user
func FindAllUserWithTodoList(c *gin.Context) {
	var user []models.User
	models.DB.Preload("TodoLists").Find(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GET '/user/todolist/:userid'
// Get user to do list by user id
func FindUserTodoListById(c *gin.Context) {
	// Get model if exist
	var user []models.User

	if err := models.DB.Where("id = ?", c.Param("id")).Preload("TodoLists").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

/*CREATE*/
//POST "/user"

// Create new userfunc
func CreateUser(c *gin.Context) {
	// Validate input
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users := models.User{UserName: input.UserName}

	models.DB.Create(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})

}
