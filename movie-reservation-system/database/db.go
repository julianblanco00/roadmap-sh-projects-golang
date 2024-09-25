package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

var Db *DB

func GetDB() (*DB, error) {
	return Db, nil
}

func Connect() {
	fmt.Println("Connecting to the database")

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		panic(err)
	}

	dbConnection := db.Ping()
	if dbConnection != nil {
		panic(dbConnection.Error())
	}

	fmt.Println("Connected to the database")

	Db = &DB{db}
}
