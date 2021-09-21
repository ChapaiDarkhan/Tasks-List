package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

)

type Tasks struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Name string `gorm:"not null" form:"name" json:"name"`
	Description  string `gorm:"not null" form:"description" json:"description"`
	Author  string `gorm:"not null" form:"author" json:"author"`
}


func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.POST("/tasks", PostTask)
		v1.GET("/tasks", GetTasks)
		v1.GET("/tasks/:id", GetTask)
		v1.PUT("/tasks/:id", UpdateTask)
		v1.DELETE("/tasks/:id", DeleteTask)
	}

	r.Run(":8080")
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	db.LogMode(true)
	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Tasks{}) {
		db.CreateTable(&Tasks{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tasks{})
	}

	return db
}

// PostTask For create new task
func PostTask(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var task Tasks
	c.Bind(&task)

	if task.Name != "" && task.Description != "" && task.Author != "" {
		// INSERT INTO "tasks" (name) VALUES (task.Name);
		db.Create(&task)
		// Display error
		c.JSON(201, gin.H{"success": task})
	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Task 1\", \"description\": \"Test task 1\", \"author\": \"user1\" }" http://localhost:8080/api/v1/tasks
}

// GetTasks For get all tasks
func GetTasks(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var tasks []Tasks
	// SELECT * FROM tasks
	db.Find(&tasks)

	// Display JSON result
	c.JSON(200, tasks)

	// curl -i http://localhost:8080/api/v1/tasks
}

// GetTask For get task by id
func GetTask(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var task Tasks
	// SELECT * FROM tasks WHERE id = *;
	db.First(&task, id)

	if task.Id != 0 {
		// Display JSON result
		c.JSON(200, task)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Task not found"})
	}

	// curl -i http://localhost:8080/api/v1/tasks/1
}

// UpdateTask For update data about task
func UpdateTask(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id task
	id := c.Params.ByName("id")
	var task Tasks
	// SELECT * FROM tasks WHERE id = *;
	db.First(&task, id)

	if task.Name != "" && task.Description != "" && task.Author != "" {

		if task.Id != 0 {
			var newTask Tasks
			c.Bind(&newTask)

			result := Tasks{
				Id:        task.Id,
				Name: newTask.Name,
				Description:  newTask.Description,
				Author:  newTask.Author,
			}

			// UPDATE tasks SET name='newTask.Name', description='newTask.description' WHERE id = task.Id;
			db.Save(&result)
			// Display modified data in JSON message "success"
			c.JSON(200, gin.H{"success": result})
		} else {
			// Display JSON error
			c.JSON(404, gin.H{"error": "Task not found"})
		}

	} else {
		// Display JSON error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"name\": \"SecondTask\", \"description\": \"Test task2\", \"author\": \"user2\" }" http://localhost:8080/api/v1/tasks/1
}

// DeleteTask For remove task
func DeleteTask(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id task
	id := c.Params.ByName("id")
	var task Tasks
	// SELECT * FROM tasks WHERE id = *;
	db.First(&task, id)

	if task.Id != 0 {
		// DELETE FROM tasks WHERE id = task.Id
		db.Delete(&task)
		// Display JSON result
		c.JSON(200, gin.H{"success": "Task #" + id + " deleted"})
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Task not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/tasks/1
}
