package main

import (
	// "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

type Todo struct {
    ID     string  `json:"id"`
    Text  string  `json:"text"`
    Completed bool  `json:"completed"`
	// Date  float64 `json:"date"` // TODO: implement latter
}
var TODOS = []Todo{
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
    var newTodo Todo

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
	var updateTodo Todo

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

func healthAndReadinessHandler(c *gin.Context){
	c.IndentedJSON(http.StatusOK, gin.H {"message": "OK"})
}

func main() {
    router := gin.Default()
    router.GET("/todos", getTodos)
	router.POST("/todos", postTodos)
	router.PUT("/todos/:id", updateTodos)
	router.DELETE("/todos/:id", deleteTodos)
	
    router.GET("/health", healthAndReadinessHandler)
    router.GET("/readiness", healthAndReadinessHandler)

	router.Run(":8080")
}
