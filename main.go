package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// https://habrahabr.ru/post/197718/
var dbHost = flag.String("db-host", "localhost", "Database hostname")
var dbPort = flag.Int("db-port", 5432, "Database port")
var dbUser = flag.String("db-user", "postgres", "Database user")
var dbPass = flag.String("db-pass", "postgres", "Database user's password")
var dbName = flag.String("db-dbname", "template0", "Database name")

//var port = flag.Int("port", "127.0.0.1", "http port")
//var host = flag.Int("host", 7788, "http host")
//var host = flag.Int("host", 7788, "http host")

// func InitDB(dataSourceName string) {
//     var err error
//     db, err = sql.Open("postgres", dataSourceName)
//     if err != nil {
//         log.Panic(err)
//     }

//     if err = db.Ping(); err != nil {
//         log.Panic(err)
//     }
// }

func main() {

	flag.Parse()

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName)

	db, err := sql.Open("postgres", connString)
	defer db.Close()

	if err != nil {
		fmt.Printf("Database opening error -->%v\n", err)
		panic("Database error")
	}

	e := echo.New()

	// Middleware
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// meta
	e.GET("/", func(c echo.Context) error {
		var resp struct {
			Collections []string `json:"collections"`
		}
		rows, err := db.Query(
			"SELECT table_name FROM information_schema.tables " +
				" WHERE table_schema = 'public'")
		if err != nil {
			return err
		}
		defer rows.Close()
		var tName string
		for rows.Next() {
			err := rows.Scan(&tName)
			if err != nil {
				return err
			}
			resp.Collections = append(resp.Collections, tName)
		}

		return c.JSON(http.StatusOK, resp)
	})

	// collection level
	e.GET("/:collection/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "GET-collection\n")
	})
	e.POST("/:collection/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "POST-collection\n")
	})

	// entity level
	e.GET("/:collection/:id/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "GET-entity\n")
	})
	e.PUT("/:collection/:id/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "UPDATE-entity\n")
	})
	e.DELETE("/:collection/:id/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "DELETE-entity\n")
	})

	// Start server
	e.Run(standard.New(":7788"))
}
