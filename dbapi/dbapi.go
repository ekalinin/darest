package dbapi

import (
	"database/sql"
	"fmt"
)

// DBApi main struct for Database Api REST
type DBApi struct {
	db *sql.DB
}

// NewPostgres open db connection for a postgres database
func NewPostgres(connstring string) (*DBApi, error) {
	return New("postgres", connstring)
}

// New open db connection for a certain dabase of a custom type
func New(dbtype string, connString string) (*DBApi, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		fmt.Printf("Database opening error: %v\n", err)
		return nil, err
	}

	return &DBApi{db}, nil
}

// GetTables returns available tables
func (d *DBApi) GetTables() ([]map[string]interface{}, error) {
	return select2map(d.db,
		"SELECT table_name as name, table_type as type "+
			"FROM information_schema.tables "+
			" WHERE table_schema = 'public'")
}

// GetTableMeta returns meta info for a table
func (d *DBApi) GetTableMeta(tableName string) ([]map[string]interface{}, error) {
	return select2map(d.db, "select column_name as name, "+
		" ordinal_position as position, column_default as default, "+
		" is_nullable as nullable, data_type as type, "+
		" is_updatable as updatable "+
		"from information_schema.columns "+
		" where table_name = '"+tableName+"'")
}

// GetEntities returns data from a table
func (d *DBApi) GetEntities(tableName string) ([]map[string]interface{}, error) {
	return select2map(d.db, "select * from "+tableName)
}

// GetEntity returns data from a table
func (d *DBApi) GetEntity(tableName string, id string) ([]map[string]interface{}, error) {
	return select2map(d.db, "select * from "+tableName+"where id="+id)
}

// Close call db.close
func (d *DBApi) Close() {
	d.db.Close()
}

//
// helpers
//

func select2map(db *sql.DB, query string) ([]map[string]interface{}, error) {
	var tableData []map[string]interface{}

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
