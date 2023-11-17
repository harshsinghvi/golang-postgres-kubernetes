package controllers

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-pg/pg/v9"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
	"time"

	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
)

// var (connection *pg.DB )// = *database.GetDatabase()

// init() {
// 	connection =  = *database.GetDatabase()
// }

// connection := *database.GetDatabase()

func GetAllTodos(c *gin.Context) {

	var todos []models.Todo
	err := database.Connection.Model(&todos).Select()
	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
}

func GetSingleTodo(c *gin.Context) {
	todoId := c.Param("id")
	todo := &models.Todo{ID: todoId}
	err := database.Connection.Select(todo)
	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    todo,
	})
}

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	c.BindJSON(&todo)

	text := todo.Text
	id := guuid.New().String()

	insertError := database.Connection.Insert(&models.Todo{
		ID:        id,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
}

func EditTodo(c *gin.Context) {
	todoId := c.Param("id")
	var todo models.Todo
	c.BindJSON(&todo)

	querry := database.Connection.Model(&models.Todo{}).Set("completed = ?", todo.Completed)

	if todo.Text != "" {
		querry.Set("text = ?", todo.Text)
	}

	res, err := querry.Where("id = ?", todoId).Update()

	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	if res.RowsAffected() == 0 {
		log.Printf("Error while update todo, Reason: \n")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo Edited Successfully",
	})
}

func DeleteTodo(c *gin.Context) {
	todoId := c.Param("id")
	todo := &models.Todo{ID: todoId}
	err := database.Connection.Delete(todo)
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
}
