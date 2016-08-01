package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strconv"

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

func select2map(db *sql.DB, query string) ([]map[string]interface{}, error) {
	tableData := make([]map[string]interface{}, 0)

	rows, err := db.Query(query)
	if err != nil {
		return tableData, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return tableData, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil

}

func main() {

	flag.Parse()

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName)

	db, err := sql.Open("postgres", connString)
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
		rows, err := select2map(db,
			"SELECT table_name as name, table_type as type "+
				"FROM information_schema.tables "+
				" WHERE table_schema = 'public'")
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, rows)
	})

	// collection level
	e.GET("/:collection/", func(c echo.Context) error {
		rows, err := select2map(db, "select * from "+c.Param("collection"))
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
		resp.Columns, err = select2map(db, "select column_name as name, "+
			"ordinal_position as position, column_default as default, "+
			"is_nullable as nullable, data_type as type, "+
			"is_updatable as updatable "+
			"from information_schema.columns "+
			" where table_name = '"+c.Param("collection")+"'")
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
		rows, err := select2map(db, "select * from "+c.Param("collection")+
			"where id="+c.Param("id"))
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
