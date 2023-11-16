package main

import (
	// "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
	models "harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"log"
	"github.com/joho/godotenv"
)


var TODOS = []models.Todo{
    {ID: "1", Text: "Task 1", Completed: false},
    {ID: "2", Text: "Task 2", Completed: false},
    {ID: "3", Text: "Task 3", Completed: false},
}

func getTodos(c *gin.Context) {
	id := c.Query("id")
	// completed := c.Query("completed") == "true" // TODO: implement this filter

	if id == "" {
		c.IndentedJSON(http.StatusOK, TODOS)
		return
	}
	
	for _, a := range TODOS {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }

	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "not found"})

}

func postTodos(c *gin.Context) {
    var newTodo models.Todo

    // Call BindJSON to bind the received JSON to
    if err := c.BindJSON(&newTodo); err != nil {
        return
    }

	if newTodo.Text == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no text"})
		return
	}

	for _, a := range TODOS {
		if(a.ID == newTodo.ID){
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Duplicate ID"})
			return
		}
	}
    // Add the new album to the slice.
    TODOS = append(TODOS, newTodo)
    c.IndentedJSON(http.StatusCreated, newTodo)
}

func updateTodos(c *gin.Context){
	id := c.Param("id")
	var updateTodo models.Todo

    // Call BindJSON to bind the received JSON to
    if err := c.BindJSON(&updateTodo); err != nil {
        return
    }
	
	if updateTodo.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no ID"})
		return
	}

	for index, a := range TODOS {
        if a.ID == id {
			if updateTodo.Text != "" {
				TODOS[index].Text = updateTodo.Text
			}
			TODOS[index].Completed = updateTodo.Completed
            c.IndentedJSON(http.StatusOK, TODOS[index])
            return
        }
    }

	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
}

func deleteTodos(c *gin.Context){
	id := c.Param("id")

	for index, a := range TODOS {
		if a.ID == id{
			TODOS = append(TODOS[:index], TODOS[index+1:]...)
			c.IndentedJSON(http.StatusOK, TODOS)
			return
		}
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no ID"})
}

func healthHandler(c *gin.Context){
	c.IndentedJSON(http.StatusOK, gin.H {"message": "OK"})
}
func readinessHandler(c *gin.Context){
	if !database.IsDtabaseReady() {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H {"message": "server not ready"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H {"message": "OK"})
}

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	
	database.Connect();
	database.CreateTodoTable();

    router := gin.Default()
    router.GET("/v1/todo", getTodos)
	router.POST("v1/todo", postTodos)
	router.PUT("v1/todo/:id", updateTodos)
	router.DELETE("v1/todo/:id", deleteTodos)

	router.GET("/v2/todo", database.GetAllTodos)
	router.GET("/v2/todos", database.GetAllTodos)
	router.GET("/v2/todo/:id", database.GetSingleTodo)
	router.POST("v2/todo", database.CreateTodo)
	router.PUT("v2/todo/:id", database.EditTodo)
	router.DELETE("v2/todo/:id", database.DeleteTodo)
	
    router.GET("/health", healthHandler)
    router.GET("/readiness", readinessHandler)

	router.Run(":8080")
}
