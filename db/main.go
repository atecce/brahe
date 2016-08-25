package db

import (
	"database/sql"
	"log"
	"reflect"
)

func constructQuery(table string, columns []string) string {

	// insert columns
	query := `INSERT INTO ` + table + `(`
	for _, column := range columns {
		query += column + ", "
	}
	query = query[0 : len(query)-2]

	// insert values
	query += `) VALUES (`
	for _ = range columns {
		query += "?, "
	}
	query = query[0 : len(query)-2]

	// update entire entry on duplicate
	query += ") ON DUPLICATE KEY UPDATE "
	for _, column := range columns {
		query += column + "=VALUES(" + column + "), "
	}
	query = query[0 : len(query)-2]

	return query
}

func addColumn(column, table string, columnType reflect.Type, canvas *sql.DB) {

	log.Printf("%s %s", column, columnType)

	// map for conversion between go and mysql data types
	goToMySQL := map[string]string{
		"bool":    "BOOL",
		"float64": "FLOAT",
		"string":  "TEXT",
	}

	// add column name and type
	query := `ALTER TABLE ` + table + ` ADD ` + column + ` ` + goToMySQL[columnType.String()]
	if result, err := canvas.Exec(query); err != nil {
		panic(err)
	} else {
		logResult(query, result)
	}
}

func AddTable(name string, canvas *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS ` + name + ` (
		id INTEGER NOT NULL,
		PRIMARY KEY (id))`
	if result, err := canvas.Exec(query); err != nil {
		panic(err)
	} else {
		logResult(query, result)
	}
}
