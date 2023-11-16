package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/controllers"
	"harshsinghvi/golang-postgres-kubernetes/controllers_old"
)

func healthHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}
func readinessHandler(c *gin.Context) {
	if !database.IsDtabaseReady() {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "server not ready"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
	
	database.Connect()
	database.CreateTodoTable()

	router := gin.Default()
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/todo", controllers_old.GetTodos)
			v1.POST("/todo", controllers_old.PostTodos)
			v1.PUT("/todo/:id", controllers_old.UpdateTodos)
			v1.DELETE("/todo/:id", controllers_old.DeleteTodos)
		}

		v2 := api.Group("/v2")
		{
			v2.GET("/todo", controllers.GetAllTodos)
			v2.GET("/todo/:id", controllers.GetSingleTodo)
			v2.POST("/todo", controllers.CreateTodo)
			v2.PUT("/todo/:id", controllers.EditTodo)
			v2.DELETE("/todo/:id", controllers.DeleteTodo)
		}
	}

	router.GET("/health", healthHandler)
	router.GET("/readiness", readinessHandler)
	router.Run(":8080")
}
