package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekalinin/darest/dbapi"
	_ "github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var dbHost = flag.String("db-host", "localhost", "Database hostname")
var dbPort = flag.Int("db-port", 5432, "Database port")
var dbUser = flag.String("db-user", "postgres", "Database user")
var dbPass = flag.String("db-pass", "postgres", "Database user's password")
var dbName = flag.String("db-dbname", "template0", "Database name")

var port = flag.Int("port", 7788, "Public http port")

func main() {

	flag.Parse()

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName)

	db, err := dbapi.NewPostgres(connString)
	defer db.Close()

	if err != nil {
		fmt.Printf("Database opening error: %v\n", err)
		panic("Database error")
	}

	e := echo.New()

	// Middleware
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// meta
	e.OPTIONS("/", func(c echo.Context) error {
		rows, err := db.GetTables()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, rows)
	})

	// collection level
	e.GET("/:collection/", func(c echo.Context) error {
		rows, err := dbapi.GetEntities(c.Param("collection"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, rows)
	})
	e.OPTIONS("/:collection/", func(c echo.Context) error {
		var resp struct {
			Pkey    []string                 `json:"pkey"`
			Columns []map[string]interface{} `json:"columns"`
		}
		// TODO: select column list for pk
		resp.Pkey = []string{"id"}
		resp.Columns, err = db.GetTableMeta(c.Param("collection"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &resp)
	})
	e.POST("/:collection/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "POST-collection\n")
	})

	// entity level
	e.GET("/:collection/:id/", func(c echo.Context) error {
		// TODO: get pk column name
		rows, err := dbapi.GetEntity(c.Param("collection"), c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, rows)
	})
	e.PUT("/:collection/:id/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "UPDATE-entity\n")
	})
	e.DELETE("/:collection/:id/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "DELETE-entity\n")
	})

	println("Server started on port: " + strconv.Itoa(*port))

	// Start server
	e.Run(standard.New(":" + strconv.Itoa(*port)))
}
