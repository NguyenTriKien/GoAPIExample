package models

import (
	"time"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TodoList struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Year     int    `json:"-"`
	Month    int    `json:"-"`
	Day      int    `json:"-"`
	CreateAt time.Time
}

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&TodoList{})
	if err != nil {
		return
	}

	DB = database
}
