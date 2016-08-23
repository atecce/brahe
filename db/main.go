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
	if result, err := canvas.Exec(`ALTER TABLE ` + table + ` ADD ` +
		column + ` ` + goToMySQL[columnType.String()]); err != nil {
		panic(err)
	} else {
		logResult(result)
	}
}

func addTable(name string, canvas *sql.DB) {
	if result, err := canvas.Exec(`CREATE TABLE IF NOT EXISTS ` + name + ` (
		id INTEGER NOT NULL,
		PRIMARY KEY (id))`); err != nil {
		panic(err)
	} else {
		logResult(result)
	}
}
