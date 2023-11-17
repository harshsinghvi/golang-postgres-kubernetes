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

	_, err := Connection.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS index_todo ON todos(id, completed, created_at, updated_at);`)

	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("Todo table and indexes created")
	return nil
}
