package main

import (
	"database/sql"
	"log"

	"repository_class/cmd/server/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// Open database connection.
	databaseConfig := &mysql.Config{
		User:      "root",
		Passwd:    "",
		Addr:      "localhost:3306",
		DBName:    "my_db",
		ParseTime: true,
	}

	db, err := sql.Open("mysql", databaseConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	// Ping database connection.
	if err = db.Ping(); err != nil {
		panic(err)
	}

	log.Println("Connection stablished")

	eng := gin.Default()
	router := routes.NewRouter(eng, db)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		panic(err)
	}
}
