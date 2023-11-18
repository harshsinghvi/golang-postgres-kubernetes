package middlewares

import (
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"harshsinghvi/golang-postgres-kubernetes/database"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"os"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
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
			return
		}

		querry := database.Connection.Model(&accessToken).Where("token = ?", token)

		if count, err = querry.Count(); err != nil {
			utils.InternalServerError(c, "Error while getting tokens, Reason:", err)
			c.Abort()
			return
		}

		if count == 0 {
			utils.UnauthorizedResponse(c)
			return
		}

		if err = querry.Select(); err != nil {
			utils.InternalServerError(c, "Error while getting all todos, Reason:", err)
			c.Abort()
			return
		}

		if time.Until(accessToken.Expiry).Seconds() <= 0 {
			utils.UnauthorizedResponse(c)
			return
		}

		c.Next()

		var hostname string
		if hostname, err = os.Hostname(); err != nil {
			log.Printf("Error loading system hostname %v\n", err)
		}
		insertError := database.Connection.Insert(&models.AccessLog{
			ID:             reqId,
			Token:          accessToken.Token,
			Path:           c.Request.URL.Path,
			ServerHostname: hostname,
			ResponseSize:   c.Writer.Size(),
			StatusCode:     c.Writer.Status(),
			ClientIP:       c.ClientIP(),
			Method:         c.Request.Method,
			ResponseTime:   time.Since(reqStart).Milliseconds(),
			CreatedAt:      time.Now(),
		})
		if insertError != nil {
			log.Println("Error loging request in db.")
			return
		}
	}
}
