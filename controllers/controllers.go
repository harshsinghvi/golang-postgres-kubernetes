package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/models/roles"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"net/http"
	"time"
)

func GetAllTodos(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var pag models.Pagination
	var err error

	var todos []models.Todo
	var searchString = c.Query("search")
	var pageString = c.Query("page")
	pag.ParseString(pageString)

	querry := database.Connection.Model(&todos).Order("created_at DESC").Where("deleted = ?", false).Where("user_id = ?", userId)

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
	userId, _ := c.Get("user_id")
	todoId := c.Param("id")
	var todos []models.Todo
	querry := database.Connection.Model(&todos).Where("id = ?", todoId).Where("deleted = ?", false).Where("user_id = ?", userId)
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
	userId, _ := c.Get("user_id")
	var todo models.Todo
	c.BindJSON(&todo)

	text := todo.Text
	id := guuid.New().String()

	newTodo := models.Todo{
		ID:        id,
		Text:      text,
		Completed: false,
		UserID:    userId.(string),
		Deleted:   false,
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
	userId, _ := c.Get("user_id")
	todoId := c.Param("id")
	var todo models.Todo
	c.BindJSON(&todo)

	querry := database.Connection.Model(&models.Todo{}).Set("completed = ?", todo.Completed).Set("updated_at = ?", time.Now()).Where("deleted = ?", false).Where("user_id = ?", userId)
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
	userId, _ := c.Get("user_id")

	querry := database.Connection.Model(&models.Todo{})

	querry = querry.Where("id = ?", todoId)
	querry = querry.Where("user_id = ?", userId)
	querry = querry.Where("deleted = ?", false)
	querry = querry.Set("deleted = ?", true)
	querry = querry.Set("updated_at = ?", time.Now())

	res, err := querry.Update()
	if err != nil {
		utils.InternalServerError(c, "error deleting todo", err)
	}
	if res.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "todo not found",
		})
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

	querry := database.Connection.Model(&models.User{}).Where("id = ?", userId).Where("deleted = ?", false)

	if count, _ := querry.Count(); count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "user_id, not found",
		})
		return
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

	querry := database.Connection.Model(&models.User{}).Where("id = ?", userId).Where("deleted = ?", false)

	if count, _ := querry.Count(); count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "user_id, not found",
		})
		return
	}

	var pag models.Pagination
	var err error
	var accessTokens []models.AccessToken
	var pageString = c.Query("page")
	pag.ParseString(pageString)

	querry = database.Connection.Model(&accessTokens).Order("created_at DESC").Where("deleted = ?", false)

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
	userId = c.Param("user-id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	querry := database.Connection.Model(&models.User{}).Where("id = ?", userId).Where("deleted = ?", false)

	if count, _ := querry.Count(); count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "user_id, not found",
		})
		return
	}

	var accessToken models.AccessToken

	c.Bind(&accessToken)

	if accessToken.Roles == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Token roles not include to update data in req body",
		})
	}

	if c.Param("user-id") == "" {
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

	querry = database.Connection.Model(&models.AccessToken{}).Set("roles = ?", accessToken.Roles)
	querry = querry.Where("id = ?", tokenId)
	querry = querry.Where("deleted = ?", false)
	querry = querry.Set("updated_at = ?", time.Now())

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
	querry := database.Connection.Model(&models.User{}).Where("id = ?", userId).Where("deleted = ?", false)

	if count, _ := querry.Count(); count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "user_id, not found",
		})
		return
	}

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
	var tokensIdStr string
	for index, accessToken := range accessTokens {
		if index == 0 {
			tokensIdStr = fmt.Sprintf("'%s'", accessToken.ID)
		}
		tokensIdStr = fmt.Sprintf("%s,'%s'", tokensIdStr, accessToken.ID)
	}

	usageQuerry := database.Connection.Model(&accessLogs)
	usageQuerry = usageQuerry.Where(fmt.Sprintf("token_id in (%s)", tokensIdStr))
	usageQuerry = usageQuerry.Where("status_code between 100 and 499")
	usageQuerry = usageQuerry.Where("billed = false")

	count, err = usageQuerry.Count()

	if err != nil {
		utils.InternalServerError(c, "Error counting users for bill, Reason:", err)
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "No Usage for Bill",
		})
		c.Abort()
		return
	}

	var bill models.Bill
	billQuerry := database.Connection.Model(&bill).Where("user_id = ?", userId)
	billQuerry = billQuerry.Where("sattled = false")
	billQuerry = billQuerry.Order("created_at DESC")
	billQuerry = billQuerry.Limit(1)
	billQuerry.Select()

	if bill.ID == "" {
		bill = models.Bill{
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
	}

	usageQuerry = usageQuerry.Set("bill_id = ?", bill.ID)
	usageQuerry = usageQuerry.Set("billed = true")
	usageQuerry = usageQuerry.Set("updated_at = ?", time.Now())
	res, updateErr := usageQuerry.Update()

	if updateErr != nil {
		log.Printf("Error While fetching access logs %s", updateErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error While fetching access logs",
		})
		c.Abort()
		return
	}

	bill.CalculateBillValue(res.RowsAffected() + bill.APIUsage)
	bill.UpdatedAt = time.Now()
	updateBillquerry := database.Connection.Model(&bill).WherePK()
	res, updateErr = updateBillquerry.Update()

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
		"status":    http.StatusOK,
		"message":   "All Bills",
		"data":      bills,
		"dueAmount": total,
	})
}

func DeleteToken(c *gin.Context) {
	var userId string
	tokenId := c.Param("token-id")
	userId = c.Param("user-id")

	if userId == "" {
		userIdFromToken, _ := c.Get("user_id")
		userId = userIdFromToken.(string)
	}

	var accessToken []models.AccessToken
	querry := database.Connection.Model(&accessToken)

	if userId != "admin" {
		querry = querry.Where("user_id = ?", userId)
	}

	querry = querry.Where("id = ?", tokenId)
	querry = querry.Where("deleted = ?", false)
	querry = querry.Set("deleted = ?", true)
	querry = querry.Set("updated_at = ?", time.Now())

	res, err := querry.Update()
	if err != nil {
		utils.InternalServerError(c, "error deleting token", err)
	}
	if res.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "token not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "token deleted",
	})
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("user-id")

	querry := database.Connection.Model(&models.User{})
	querry = querry.Where("id = ?", userId)
	querry = querry.Where("deleted = ?", false)
	querry = querry.Set("deleted = ?", true)
	querry = querry.Set("updated_at = ?", time.Now())

	res, err := querry.Update()
	if err != nil {
		utils.InternalServerError(c, "error deleting user", err)
		return
	}

	if res.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	querry = database.Connection.Model(&models.AccessToken{})
	querry = querry.Where("user_id = ?", userId)
	querry = querry.Where("deleted = ?", false)
	querry = querry.Set("deleted = ?", true)
	querry = querry.Set("updated_at = ?", time.Now())

	res, err = querry.Update()
	if err != nil {
		utils.InternalServerError(c, "error deleting user", err)
		return
	}
	count := res.RowsAffected()
	if count == 0 {
		log.Printf("error while deleting tokens of user %s", userId)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"message":       "user and its token deleted",
		"tokensDeleted": count,
	})
}
