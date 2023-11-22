package view

import (
	"github.com/gin-gonic/gin"

	"Module/API/models"

	"Module/API/controllers"
)

func Response() {
	r := gin.Default()

	models.ConnectDatabase() // new

	r.GET("/list", controllers.FindAll)
	r.GET("/list/:id", controllers.FindListById)
	r.POST("/list", controllers.CreateToDo)
	r.PUT("/list/:id", controllers.UpdateToDoList)
	r.DELETE("/list/:id", controllers.DeleteListItem)

	err := r.Run()
	if err != nil {
		return
	}
}
