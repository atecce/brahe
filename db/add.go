package db

import (
	"bodhi/herodotus"
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var dbLog = herodotus.CreateFileLog("db")

func Initiate() *sql.DB {

	// connect
	canvas, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")

	// create
	query := `CREATE DATABASE IF NOT EXISTS canvas`
	_, err = canvas.Exec(query)
	dbLog.Println(query)

	// use
	query = `USE canvas`
	_, err = canvas.Exec(query)
	dbLog.Println(query)

	if err != nil {
		panic(err)
	}

	return canvas
}

func AddTable(name string, canvas *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS ` + name + ` (id INTEGER NOT NULL, PRIMARY KEY (id))`
	result, err := canvas.Exec(query)
	dbLog.Println(query, result)

	if err != nil {
		panic(err)
	}
}

func addColumn(column, table string, columnType reflect.Type, canvas *sql.DB) {

	// map for conversion between go and mysql data types
	goToMySQL := map[string]string{
		"bool":    "BOOL",
		"float64": "FLOAT",
		"string":  "TEXT",
	}

	// add column name and type
	query := `ALTER TABLE ` + table + ` ADD ` + column + ` ` + goToMySQL[columnType.String()]
	_, err := canvas.Exec(query)
	dbLog.Println(query)

	if err != nil {
		panic(err)
	}
}

func splitMap(row map[string]interface{}, canvas *sql.DB) ([]string, []interface{}) {

	// add columns and values to query
	var columns []string
	var values []interface{}
	for column, value := range row {

		// make sure field isn't empty
		if value != nil && value != "" {

			// create a new table for additional map
			if reflect.ValueOf(value).Kind() == reflect.Map {
				AddTable(column, canvas)

				// recursion WTF
				AddRow(column, value.(map[string]interface{}), canvas)
				continue

			} else {

				// special case for reserved MySQL word
				var entry string
				if column == "release" {
					entry = "release_number"
				} else {
					entry = column
				}

				// append columns and values to list
				columns = append(columns, entry)
				values = append(values, value)
			}
		}
	}

	return columns, values
}

func AddRow(table string, row map[string]interface{}, canvas *sql.DB) {

	// split map into columns and values
	columns, values := splitMap(row, canvas)

	// construct query out of lists
	query := constructQuery(table, columns)

	// prepare statment
	stmt, err := canvas.Prepare(query)
	if err != nil {

		// assert error is MySQL specific
		prepareErr := err.(*mysql.MySQLError)

		// handle unknown column
		if prepareErr.Number == 1054 {

			// column name is second field delimited with single quotes
			unknownColumn := strings.Split(prepareErr.Message, "'")[1]

			// handle special case for MySQL keyword
			var columnType reflect.Type
			if unknownColumn == "release_number" {
				columnType = reflect.TypeOf(row["release"])
			} else {
				columnType = reflect.TypeOf(row[unknownColumn])
			}

			// add column and try to add the row again
			addColumn(unknownColumn, table, columnType, canvas)
			AddRow(table, row, canvas)

		} else {
			panic(prepareErr)
		}
	} else {
		defer stmt.Close()

		// insert row
		_, err := stmt.Exec(values...)
		if err != nil {

			// assert error is MySQL specific
			execErr := err.(*mysql.MySQLError)

			// handle bananas characters
			if execErr.Number == 1366 {

				// error message convenienty delimted by '
				problemColumn := strings.Split(execErr.Message, "'")[3]
				problemValue := strings.Split(execErr.Message, "'")[1]

				// reset as string
				row[problemColumn] = problemValue

			} else {
				panic(execErr)
			}
		} else {

			// log only values of insert query
			dbLog.Println("INSERT INTO ", columns, "VALUES ", values)
		}
	}
}
