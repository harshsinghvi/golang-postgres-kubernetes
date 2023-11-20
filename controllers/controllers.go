package controllers

import (
	"fmt"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/models/roles"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
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
	var todos []models.Todo
	querry := database.Connection.Model(&todos).Where("id = ?", todoId)
	if count, _ := querry.Count(); count == 1 {
		if err := querry.Select(); err != nil {
			utils.InternalServerError(c, "Error while getting a single todo, Reason:", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Single Todo",
			"data":    todos,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    todos,
	})
}

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	c.BindJSON(&todo)

	text := todo.Text
	id := guuid.New().String()

	newTodo := models.Todo{
		ID:        id,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertError := database.Connection.Insert(&newTodo)

	if insertError != nil {
		utils.InternalServerError(c, "Error while inserting new todo into db, Reason:", insertError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    newTodo,
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
	var userId string
	userId = c.Param("id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	if userId != "admin" {
		if count, _ := database.Connection.Model(&models.User{}).Where("id = ?", userId).Count(); count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "user_id, not found",
			})
			return
		}
	}

	id := guuid.New().String()
	token := utils.GenerateToken(id)
	accessToken := models.AccessToken{
		ID:        id,
		Token:     token,
		UserID:    userId,
		Roles:     []string{roles.Read, roles.Write},
		Expiry:    time.Now().AddDate(0, 0, 10),
		CreatedAt: time.Now(),
	}
	insertError := database.Connection.Insert(&accessToken)

	if insertError != nil {
		utils.InternalServerError(c, "Error while inserting new token into db, Reason:", insertError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Token created Successfully",
		"data":    accessToken,
		"token":   token,
	})
}

func GetTokens(c *gin.Context) {
	var userId string
	userId = c.Param("id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	var pag models.Pagination
	var err error
	var accessTokens []models.AccessToken
	var pageString = c.Query("page")
	pag.ParseString(pageString)

	querry := database.Connection.Model(&accessTokens).Order("created_at DESC")

	if userId != "admin" {
		querry = querry.Where("user_id = ?", userId)
	}

	if pag.TotalRecords, err = querry.Count(); err != nil {
		utils.InternalServerError(c, "Error while getting tokens, Reason:", err)
		return
	}

	if pag.CurrentPage != -1 {
		querry = querry.Limit(10).Offset(10 * (pag.CurrentPage))
	}

	if err = querry.Select(); err != nil {
		utils.InternalServerError(c, "Error while getting Tokens, Reason:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    fmt.Sprintf("All Tokens by %s", userId),
		"data":       accessTokens,
		"pagination": pag.Validate(),
	})
}

func UpdateToken(c *gin.Context) {

	var userId string
	tokenId := c.Param("token-id")
	userId = c.Param("id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	var accessToken models.AccessToken

	c.Bind(&accessToken)

	if accessToken.Roles == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Token roles not include to update data in req body",
		})
	}

	if c.Param("id") == "" {
		for _, role := range accessToken.Roles {
			if role == roles.Admin {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "invalid role admin",
				})
				return
			}
		}
	}

	querry := database.Connection.Model(&models.AccessToken{}).Set("roles = ?", accessToken.Roles).Set("updated_at = ?", time.Now())
	querry = querry.Where("id = ?", tokenId)

	if userId != "admin" {
		querry = querry.Where("user_id = ?", userId)
	}

	res, err := querry.Update()
	if err != nil {
		utils.InternalServerError(c, "Error while editing token, Reason:", err)
	}

	if res.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Token/user not found or unauthorised request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Token Edited Successfully",
	})
}

func CreateNewUser(c *gin.Context) {
	userId := guuid.New().String()
	tokenId := guuid.New().String()
	token := utils.GenerateToken(tokenId)

	user := models.User{}

	c.Bind(&user)

	count, err := database.Connection.Model(&models.User{}).Where("email = ?", user.Email).Count()
	log.Println(count)

	if err != nil {
		utils.InternalServerError(c, "Error while getting tokens, Reason:", err)
		return
	}

	if count != 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Email already exists",
		})
		return
	}

	user = models.User{
		ID:        userId,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}

	accessToken := models.AccessToken{
		ID:        tokenId,
		Token:     token,
		UserID:    userId,
		Roles:     []string{roles.Read, roles.Write},
		Expiry:    time.Now().AddDate(0, 0, 10),
		CreatedAt: time.Now(),
	}

	if insertError := database.Connection.Insert(&user); insertError != nil {
		utils.InternalServerError(c, "Error while inserting new user into db, Reason:", insertError)
		return
	}

	if insertError := database.Connection.Insert(&accessToken); insertError != nil {
		utils.InternalServerError(c, "Error while inserting new token into db, Reason:", insertError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User And Token created Successfully",
		"user":    user,
		"token":   accessToken,
	})
}

func GetUserID(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var accessTokens []models.AccessToken
	database.Connection.Model(&accessTokens).Where("user_id = ?", userId.(string)).Order("created_at DESC").Select()

	c.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"user_id":       userId,
		"access_tokens": accessTokens,
	})
}

func CreateBill(c *gin.Context) {
	var userId string
	userId = c.Param("id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	var err error

	count, err := database.Connection.Model(&models.User{}).Where("id = ?", userId).Count()
	if err != nil {
		utils.InternalServerError(c, "Error counting users for bill, Reason:", err)
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "user not found",
		})
		c.Abort()
		return
	}

	var bill = models.Bill{
		ID:        guuid.New().String(),
		APIUsage:  0,
		BillValue: 0,
		Sattled:   false,
		UserID:    userId,
		CreatedAt: time.Now(),
	}

	if insertError := database.Connection.Insert(&bill); insertError != nil {
		utils.InternalServerError(c, "Error inserting bill, Reason:", insertError)
		return
	}

	var accessTokens []models.AccessToken

	if err = database.Connection.Model(&accessTokens).Where("user_id = ?", userId).Select(); err != nil {
		log.Printf("no Token found for the given user %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "no Token found for the given user",
		})
		c.Abort()
		return
	}

	var accessLogs []models.AccessLog

	var tokensStr string

	for index, accessToken := range accessTokens {
		if index == 0 {
			tokensStr = fmt.Sprintf("'%s'", accessToken.Token)
		}
		tokensStr = fmt.Sprintf("%s,'%s'", tokensStr, accessToken.Token)
	}

	querry := database.Connection.Model(&accessLogs)
	querry = querry.Set("bill_id = ?", bill.ID)
	querry = querry.Set("billed = true")

	querry = querry.Where(fmt.Sprintf("token in (%s)", tokensStr))
	querry = querry.Where("status_code between 100 and 499")
	querry = querry.Where("billed = false")

	res, updateErr := querry.Update()

	if updateErr != nil {
		log.Printf("Error While fetching access logs %s", updateErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error While fetching access logs",
		})
		c.Abort()
		return
	}

	// if res.RowsAffected() == 0 {
	// TODO: Delete bill and return error
	// }

	bill.CalculateBillValue(res.RowsAffected())

	querry = database.Connection.Model(&bill).WherePK()
	res, updateErr = querry.Update()

	if updateErr != nil || res.RowsAffected() == 0 {
		log.Printf("Error While updating bill %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error While updating bill",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Billing Done",
		"data":    bill,
	})
}

func GetBills(c *gin.Context) {
	var userId string
	userId = c.Param("id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	var err error
	var total float32 = 0
	var bills []models.Bill

	err = database.Connection.Model(&bills).Where("user_id = ?", userId).Order("created_at DESC").Select()

	if err != nil {
		log.Printf("Error while getting bills, Reason: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error While getting bills or not bills found",
		})
	}

	for _, bill := range bills {
		if !bill.Sattled {
			total += bill.BillValue
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Bills",
		"data":    bills,
		"total":   total,
	})
}
