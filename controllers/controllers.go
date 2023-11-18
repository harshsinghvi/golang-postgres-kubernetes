package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"net/http"
	"time"
)

func GetAllTodos(c *gin.Context) {
	var pag models.Pagination
	var err error

	var todos []models.Todo
	var searchString = c.Query("search")
	var pageString = c.Query("page")
	pag.ParseString(pageString)

	querry := database.Connection.Model(&todos).Order("created_at DESC")

	if searchString != "" {
		querry = querry.Where(fmt.Sprintf("text like '%%%s%%'", searchString))
	}

	if pag.TotalRecords, err = querry.Count(); err != nil {
		utils.InternalServerError(c, "Error while getting all todos, Reason:", err)
		return
	}

	if pag.CurrentPage != -1 {
		querry = querry.Limit(10).Offset(10 * (pag.CurrentPage))
	}

	if err := querry.Select(); err != nil {
		utils.InternalServerError(c, "Error while getting all todos, Reason:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "All Todos",
		"data":       todos,
		"pagination": pag.Validate(),
	})
}

func GetSingleTodo(c *gin.Context) {
	todoId := c.Param("id")
	todo := &models.Todo{ID: todoId}
	if err := database.Connection.Select(todo); err != nil {
		utils.InternalServerError(c, "Error while getting a single todo, Reason:", err)
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
		utils.InternalServerError(c, "Error while inserting new todo into db, Reason:", insertError)
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

	querry := database.Connection.Model(&models.Todo{}).Set("completed = ?", todo.Completed).Set("updated_at = ?", time.Now())
	if todo.Text != "" {
		querry = querry.Set("text = ?", todo.Text)
	}

	res, err := querry.Where("id = ?", todoId).Update()

	if err != nil {
		utils.InternalServerError(c, "Error while editing todo, Reason:", err)
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
	if err := database.Connection.Delete(todo); err != nil {
		utils.InternalServerError(c, "Error while deleting a single todo, Reason:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
}

func CreateNewToken(c *gin.Context) {
	id := guuid.New().String()
	token := utils.GenerateToken(id)
	var accessToken models.AccessToken
	c.BindJSON(&accessToken)

	insertError := database.Connection.Insert(&models.AccessToken{
		ID:        id,
		Token:     token,
		Email:     accessToken.Email, // TODO: validate Email
		Expiry:    time.Now().AddDate(0, 0, 10),
		CreatedAt: time.Now(),
	})

	if insertError != nil {
		utils.InternalServerError(c, "Error while inserting new todo into db, Reason:", insertError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Token created Successfully",
		"token":   token,
	})
}

// Get tokens by email

func GetTokens(c *gin.Context) {
	email := c.Param("email")
	var pag models.Pagination
	var err error
	var accessTokens []models.AccessToken
	var pageString = c.Query("page")
	pag.ParseString(pageString)

	querry := database.Connection.Model(&accessTokens).Order("created_at DESC")

	if email != "admin" {
		querry = querry.Where("email = ?", email)
	}

	if pag.TotalRecords, err = querry.Count(); err != nil {
		utils.InternalServerError(c, "Error while getting tokens, Reason:", err)
		return
	}

	if pag.CurrentPage != -1 {
		querry = querry.Limit(10).Offset(10 * (pag.CurrentPage))
	}

	if err := querry.Select(); err != nil {
		utils.InternalServerError(c, "Error while getting Tokens, Reason:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    fmt.Sprintf("All Tokens by %s", email),
		"data":       accessTokens,
		"pagination": pag.Validate(),
	})
}
