package view

import (
	"github.com/gin-gonic/gin"

	"Module/API/models"

	"Module/API/controllers"
)

func Response() {
	r := gin.Default()

	models.ConnectDatabase() // new

	//To do list crud
	r.GET("/list", controllers.FindAll)               // find all to-do list
	r.GET("/list/:id", controllers.FindListById)      // find a list by id
	r.POST("/list", controllers.CreateToDo)           // create a list
	r.PUT("/list/:id", controllers.UpdateToDoList)    // update a to-do list base on id
	r.DELETE("/list/:id", controllers.DeleteListItem) // delete a to-do list base on id

	//User crud
	r.POST("/user", controllers.CreateUser)                       // create a user
	r.GET("/user", controllers.FindAllUser)                       // find all user
	r.GET("/user/todolist", controllers.FindAllUserWithTodoList)  // find all user with their to-do list
	r.GET("/user/todolist/:id", controllers.FindUserTodoListById) // find a user with their to-do list

	err := r.Run()
	if err != nil {
		return
	}
}
