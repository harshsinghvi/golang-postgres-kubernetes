package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"harshsinghvi/golang-postgres-kubernetes/controllers"
	"harshsinghvi/golang-postgres-kubernetes/controllers_old"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/middlewares"
	"log"
	"net/http"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	}
}
func main() {

	database.Connect()
	database.CreateTables()

	// router := gin.Default()
	// TODO: for improved server fault tollerent but dosent log requests
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
			v2.POST("/token", controllers.CreateNewToken)
			v2.GET("/token/:email", controllers.GetTokens)

			todo := v2.Group("/todo")
			{
				todo.Use(middlewares.AuthMiddleware())
				todo.GET("/", controllers.GetAllTodos)
				todo.GET("/:id", controllers.GetSingleTodo)
				todo.POST("/", controllers.CreateTodo)
				todo.PUT("/:id", controllers.EditTodo)
				todo.DELETE("/:id", controllers.DeleteTodo)
			}

		}
	}

	router.GET("/health", healthHandler)
	router.GET("/readiness", readinessHandler)
	router.Run(":8080")
}
