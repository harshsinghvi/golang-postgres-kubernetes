package middlewares

import (
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/models/roles"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"os"
	"time"
)

type Config map[string]interface{}

func parseConfigOption(config []Config, key string, defaultValue interface{}) interface{} {
	if config == nil {
		return defaultValue
	}
	if val, ok := config[0][key]; ok {
		return val
	}
	return defaultValue
}

func AIO(requiredRoles roles.Roles, config ...Config) gin.HandlerFunc {
	billingDisable := parseConfigOption(config, "billing-disable", false).(bool)

	return func(c *gin.Context) {
		reqStart := time.Now()
		var accessToken models.AccessToken
		var count int
		var err error
		token := c.GetHeader("token")
		reqId := guuid.New().String()

		c.Set("requestId", reqId)
		c.Writer.Header().Set("X-Request-Id", reqId)

		if token == "" {
			utils.UnauthorizedResponse(c)
			logReqToDb(reqId, accessToken, c, reqStart, false)
			return
		}

		querry := database.Connection.Model(&accessToken).Where("token = ?", token).Where("deleted = ?", false)

		if count, err = querry.Count(); err != nil {
			utils.InternalServerError(c, "Error while getting tokens, Reason:", err)
			c.Abort()
			logReqToDb(reqId, accessToken, c, reqStart, false)
			return
		}

		if count == 0 {
			utils.UnauthorizedResponse(c)
			logReqToDb(reqId, accessToken, c, reqStart, billingDisable)
			return
		}

		if err = querry.Select(); err != nil {
			utils.InternalServerError(c, "Error while getting all todos, Reason:", err)
			logReqToDb(reqId, accessToken, c, reqStart, false)
			return
		}

		if time.Until(accessToken.Expiry).Seconds() <= 0 ||
			!roles.CheckRoles(requiredRoles, accessToken.Roles) {
			utils.UnauthorizedResponse(c)
			logReqToDb(reqId, accessToken, c, reqStart, billingDisable)
			return
		}

		c.Set("token", token)
		c.Set("user_id", accessToken.UserID)
		c.Next()
		logReqToDb(reqId, accessToken, c, reqStart, billingDisable)
	}
}

func logReqToDb(reqId string, accessToken models.AccessToken, c *gin.Context, reqStart time.Time, billingDisable bool) {
	var err error
	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		log.Printf("Error loading system hostname %v\n", err)
	}

	insertError := database.Connection.Insert(&models.AccessLog{
		ID:             reqId,
		TokenID:        accessToken.ID,
		Path:           c.Request.URL.Path,
		ServerHostname: hostname,
		ResponseSize:   c.Writer.Size(),
		StatusCode:     c.Writer.Status(),
		ClientIP:       c.ClientIP(),
		Method:         c.Request.Method,
		ResponseTime:   time.Since(reqStart).Milliseconds(),
		CreatedAt:      time.Now(),
		Billed:         billingDisable, // Billed already true of not needed to be billed
		BillID:         "Not Needed/Billed",
	})
	if insertError != nil {
		log.Println("Error loging request in db.")
		return
	}
}
