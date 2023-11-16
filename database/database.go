package database

import (
	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	"harshsinghvi/golang-postgres-kubernetes/models"
	"harshsinghvi/golang-postgres-kubernetes/utils"
	"log"
)

var database_ready = false
var Connection *pg.DB

func IsDtabaseReady() bool {
	return database_ready
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
	database_ready = true
	return Connection
}

// Create User Table
func CreateTodoTable() error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	createError := Connection.CreateTable(&models.Todo{}, opts)
	if createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Todo table created")
	return nil
}
