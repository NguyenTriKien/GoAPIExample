package controllers

import (
	"net/http"
	"time"

	"Module/API/models"

	"github.com/gin-gonic/gin"
)

// / GET ALL ///
// GET /list
// Get all list
/* func FindAll(c *gin.Context) {
	var todoList []models.TodoList // khai báo một todoList là môt slice vì todoList cần mở rộng để lưu trữ thêm nhiều todoList khác
	models.DB.Find(&todoList)

	c.JSON(http.StatusOK, gin.H{"data": todoList})
} */

// /CREATE///
type CreateToDoListInput struct {
	Title    string    `json:"title" binding:"required"`
	Status   string    `json:"status"`
	Year     int       `json:"year" binding:"required"`
	Month    int       `json:"month" binding:"required"`
	Day      int       `json:"day" binding:"required"`
	CreateAt time.Time `json:"CreateAt"`
	UserID   uint      `json:"userid"`
}

// check valid date of a month in a specific year
func daysInMonth(year int, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// POST /list
// Create new to do work
func CreateToDo(c *gin.Context) {
	// Validate input
	var input CreateToDoListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get UserID from token
	userID, exists := getUserIDFromToken(c)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if input.Month > 12 || input.Month <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month input"})
		return
	}

	if input.Day > 31 || input.Day <= 0 || input.Day > daysInMonth(input.Year, int(input.Month)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day input"})
		return
	}

	// Create to do list
	todoList := models.TodoList{
		Title:  input.Title,
		Status: input.Status,
		CreateAt: time.Date(input.Year,
			time.Month(input.Month),
			input.Day,
			time.Now().Hour(),
			time.Now().Minute(),
			time.Now().Second(),
			time.Now().Nanosecond(),
			time.Local),
		UserID: userID}

	//Điều kiện kiểm tra thời hạn của to do list
	createdAt := todoList.CreateAt

	if createdAt.Equal(time.Now()) {
		todoList.Status = "Due"
	}
	if createdAt.Before(time.Now()) {
		todoList.Status = "Over due"
	}
	if createdAt.After(time.Now()) {
		todoList.Status = "Open"
	}

	models.DB.Create(&todoList)

	c.JSON(http.StatusOK, gin.H{"data": todoList})
}

// GET '/user/todolist/:id'
// Get user to do list by user id
func FindUserTodoListById(c *gin.Context) {
	// Get model if exist
	var user []models.User

	if err := models.DB.Where("id = ?", c.Param("userid")).Preload("TodoLists").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// /GET BY ID///
// GET /lists/:id
// Find a lists
/* func FindListById(c *gin.Context) {
	// Get model if exist
	var todoList models.TodoList

	if err := models.DB.Where("id = ?", c.Param("id")).First(&todoList).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todoList})
}
*/
///UPDATE///

type UpdateToDoListInput struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

// PATCH /list/:id
// Update a book
func UpdateToDoList(c *gin.Context) {
	// Get model if exist
	var todoList models.TodoList
	if err := models.DB.Where("id = ?", c.Param("id")).First(&todoList).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateToDoListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&todoList).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": todoList})
}

// /Delete///
// DELETE /list/:id
// Delete an item in a list
func DeleteListItem(c *gin.Context) {
	// Get model if exist
	var todoList models.TodoList
	if err := models.DB.Where("id = ?", c.Param("id")).First(&todoList).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&todoList)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
