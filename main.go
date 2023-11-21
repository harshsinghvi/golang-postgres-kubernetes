package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"harshsinghvi/golang-postgres-kubernetes/controllers"
	"harshsinghvi/golang-postgres-kubernetes/controllers_old"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/middlewares"
	"harshsinghvi/golang-postgres-kubernetes/models/roles"
	"harshsinghvi/golang-postgres-kubernetes/utils"
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
	PORT := utils.GetEnv("PORT", "8080")

	database.Connect()
	database.CreateTables()

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
			// Public endpoint
			v2.POST("/user", controllers.CreateNewUser)

			// Admin Endpoints
			v2.POST("/user/:id/token", middlewares.AIO(roles.Roles{}), controllers.CreateNewToken)
			v2.GET("/user/:id/token", middlewares.AIO(roles.Roles{}), controllers.GetTokens)
			v2.PUT("/user/:user-id/token/:token-id", middlewares.AIO(roles.Roles{}), controllers.UpdateToken)
			v2.POST("/user/:id/bill", middlewares.AIO(roles.Roles{}, middlewares.Config{"billing-disable": true}), controllers.CreateBill)
			v2.GET("/user/:id/bill", middlewares.AIO(roles.Roles{}, middlewares.Config{"billing-disable": true}), controllers.GetBills)
			v2.DELETE("/user/:user-id", middlewares.AIO(roles.Roles{}), controllers.DeleteUser)
			v2.DELETE("/user/:user-id/token/:token-id", middlewares.AIO(roles.Roles{}), controllers.DeleteToken)

			// Users Endpoints
			v2.GET("/user", middlewares.AIO(roles.Roles{roles.Any}), controllers.GetUserID)
			v2.POST("/user/token", middlewares.AIO(roles.Roles{roles.Write}), controllers.CreateNewToken)
			v2.GET("/user/token", middlewares.AIO(roles.Roles{roles.Read}), controllers.GetTokens)
			v2.PUT("/user/token/:token-id", middlewares.AIO(roles.Roles{roles.Write}), controllers.UpdateToken)
			v2.POST("/user/bill", middlewares.AIO(roles.Roles{roles.Any}, middlewares.Config{"billing-disable": true}), controllers.CreateBill)
			v2.GET("/user/bill", middlewares.AIO(roles.Roles{roles.Any}, middlewares.Config{"billing-disable": true}), controllers.GetBills)
			v2.DELETE("/user/token/:token-id", middlewares.AIO(roles.Roles{roles.Write}), controllers.DeleteToken)

			// Business Logic
			v2.GET("/todo/", middlewares.AIO(roles.Roles{roles.Admin, roles.Read}), controllers.GetAllTodos)
			v2.GET("/todo/:id", middlewares.AIO(roles.Roles{roles.Admin, roles.Read, roles.ReadOne}), controllers.GetSingleTodo)
			v2.POST("/todo/", middlewares.AIO(roles.Roles{roles.Admin, roles.Write, roles.WriteNewOnly}), controllers.CreateTodo)
			v2.PUT("/todo/:id", middlewares.AIO(roles.Roles{roles.Admin, roles.Write, roles.WriteUpdateOnly}), controllers.EditTodo)
			v2.DELETE("/todo/:id", middlewares.AIO(roles.Roles{roles.Admin, roles.Write}), controllers.DeleteTodo)
		}
	}

	router.GET("/health", healthHandler)
	router.GET("/readiness", readinessHandler)

	router.Run(fmt.Sprintf(":%s", PORT))
}
