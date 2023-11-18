package database

import (
	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
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

func CreateTables() error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	if createError := Connection.CreateTable(&models.Todo{}, opts); createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}
	if _, err := Connection.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS go_index_todos ON todos(completed, created_at);`); err != nil {
		log.Println(err.Error())
		return err
	}
	if createError := Connection.CreateTable(&models.AccessToken{}, opts); createError != nil {
		log.Printf("Error while creating access_tokens table, Reason: %v\n", createError)
		return createError
	}
	if _, err := Connection.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS go_index_access_tokens ON access_tokens(created_at, token, email);`); err != nil {
		log.Println(err.Error())
		return err
	}
	if createError := Connection.CreateTable(&models.AccessLog{}, opts); createError != nil {
		log.Printf("Error while creating access_logs table, Reason: %v\n", createError)
		return createError
	}
	// TODO
	if _, err := Connection.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS go_index_access_logs ON access_logs(token, path, method, response_time, status_code, server_hostname, created_at);`); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Todo table and indexes created")
	return nil
}
