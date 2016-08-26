package db

import (
	"database/sql"
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func Initiate() *sql.DB {

	// create database
	if canvas, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"); err != nil {
		panic(err)
	} else {
		query := `CREATE DATABASE IF NOT EXISTS canvas`
		if result, err := canvas.Exec(query); err != nil {
			panic(err)
		} else {
			logResult(query, result)

			// use the database
			query = `USE canvas`
			if result, err := canvas.Exec(query); err != nil {
				panic(err)
			} else {
				logResult(query, result)

				return canvas
			}
		}
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

func AddRow(table string, row map[string]interface{}, canvas *sql.DB) {

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

	// construct query out of lists
	query := constructQuery(table, columns)

	// log.Println(query)

	// prepare statment
	if stmt, err := canvas.Prepare(query); err != nil {
		if prepareErr, ok := err.(*mysql.MySQLError); ok {

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
				panic(err)
			}
		}
	} else {
		defer stmt.Close()

		// insert row
		if result, err := stmt.Exec(values...); err != nil {
			if execErr, ok := err.(*mysql.MySQLError); ok {

				// handle bananas characters
				if execErr.Number == 1366 {
					log.Println(values)

					// error message convenienty delimted by '
					problemColumn := strings.Split(execErr.Message, "'")[3]
					problemValue := strings.Split(execErr.Message, "'")[1]
					log.Println(problemColumn, problemValue, reflect.TypeOf(problemValue))

					// reset as string
					row[problemColumn] = problemValue
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {

			logResult(query, result)
		}
	}
}
