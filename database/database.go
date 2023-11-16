package database

import (
	. "harshsinghvi/golang-postgres-kubernetes/models"

	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	guuid "github.com/google/uuid"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var database_ready = false
var connection *pg.DB

func IsDtabaseReady() bool {
	return database_ready
}

func GetDatabase() *pg.DB {
	return connection
}

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "postgres",
		Addr:     "localhost:5432",
		Database: "postgres",
	}
	connection = pg.Connect(opts)
	if connection == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")
	return connection
}

// Create User Table
func CreateTodoTable() error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := connection.CreateTable(&Todo{}, opts)
	if createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Todo table created")
	return nil
}

func GetAllTodos(c *gin.Context) {
	var todos []Todo
	err := connection.Model(&todos).Select()
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
	todo := &Todo{ID: todoId}
	err := connection.Select(todo)
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
	var todo Todo
	c.BindJSON(&todo)

	text := todo.Text
	id := guuid.New().String()

	insertError := connection.Insert(&Todo{
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
	var todo Todo
	c.BindJSON(&todo)
	completed := todo.Completed
	_, err := connection.Model(&Todo{}).Set("completed = ?", completed).Where("id = ?", todoId).Update()
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
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
	todo := &Todo{ID: todoId}
	err := connection.Delete(todo)
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
