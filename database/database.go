package database

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	guuid "github.com/google/uuid"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/models/roles"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
	"time"
)

var Connection *pg.DB

func IsDtabaseReady() bool {
	ctx := Connection.Context()
	var version string
	_, err := Connection.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()")
	if err != nil {
		log.Printf("Failed to connect to database")
		return false
	}
	return true
}

func GetDatabase() **pg.DB {
	return &Connection
}

func Connect() *pg.DB {
	DB_HOST := utils.GetEnv("DB_HOST", "localhost")
	DB_PORT := utils.GetEnv("DB_PORT", "5432")
	DB_USER := utils.GetEnv("DB_USER", "postgres")
	DB_PASSWORD := utils.GetEnv("DB_PASSWORD", "postgres")
	DB_NAME := utils.GetEnv("DB_NAME", "postgres")

	opts := &pg.Options{
		User:     DB_USER,
		Password: DB_PASSWORD,
		Addr:     DB_HOST + ":" + DB_PORT,
		Database: DB_NAME,
	}

	Connection = pg.Connect(opts)

	if Connection == nil {
		log.Printf("Failed to connect to database")
		return nil
	}

	ctx := Connection.Context()
	var version string
	_, err := Connection.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()")
	if err != nil {
		log.Printf("Failed to connect to database")
		return nil
	}

	log.Printf("Connected to db")

	return Connection
}
func createTablesAndIndexes(tableName string, model interface{}, indexFields string) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	createIndexQuerry := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS go_index_%s ON %s(%s);", tableName, tableName, indexFields)

	if createError := Connection.CreateTable(model, opts); createError != nil {
		log.Printf("Error while creating %s table, Reason: %v\n", tableName, createError)
		return createError
	}

	if _, err := Connection.Exec(createIndexQuerry); err != nil {
		log.Printf("Error while index of %s table, Reason: %v\n", tableName, err)
		return err
	}

	log.Printf("INFO: %s table and its indexes created", tableName)
	return nil
}

func CreateTables() {
	createTablesAndIndexes("todos", &models.Todo{}, "completed, created_at, deleted")
	createTablesAndIndexes("access_tokens", &models.AccessToken{}, "created_at, token, expiry, user_id, deleted")
	createTablesAndIndexes("access_logs", &models.AccessLog{}, "token_id, path, method, response_time, status_code, server_hostname, created_at, bill_id, billed, deleted")
	createTablesAndIndexes("users", &models.User{}, "email, created_at, deleted")
	createTablesAndIndexes("bills", &models.Bill{}, "sattled, user_id, created_at, deleted")
	checkAndCreateAdminUser()
	checkAndCreateAdminToken()
	log.Printf("all tables and indexes created")
}

func checkAndCreateAdminToken() {
	var accessToken models.AccessToken
	querry := Connection.Model(&accessToken).Where("id = ?", "admin")
	count, err := querry.Count()
	if err != nil {
		log.Println("Error in getting access_token count")
	}
	if count != 0 {
		return
	}

	id := "admin"
	token := utils.GenerateToken(guuid.New().String())

	insertError := Connection.Insert(&models.AccessToken{
		ID:        id,
		Token:     token,
		UserID:    id,
		Expiry:    time.Now().AddDate(99, 0, 00),
		CreatedAt: time.Now(),
		Roles:     []string{roles.Admin},
	})

	if insertError != nil {
		log.Printf("Error while inserting new token into db, Reason: %v\n", insertError)
	}

	log.Printf("Admin Token created")
}

func checkAndCreateAdminUser() {
	var user models.User
	querry := Connection.Model(&user).Where("id = ?", "admin")
	count, err := querry.Count()
	if err != nil {
		log.Println("Error in getting users count")
	}
	if count != 0 {
		return
	}

	id := "admin"

	insertError := Connection.Insert(&models.User{
		ID:        id,
		Email:     id,
		CreatedAt: time.Now(),
	})

	if insertError != nil {
		log.Printf("Error while inserting new user into db, Reason: %v\n", insertError)
	}

	log.Printf("Admin user created")
}
