package db

import (
	"database/sql"
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func InitiateDB() *sql.DB {

	// prepare db
	if canvas, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"); err != nil {
		panic(err)
	} else {

		// create database
		if result, err := canvas.Exec(`CREATE DATABASE IF NOT EXISTS canvas`); err != nil {
			panic(err)
		} else {

			logResult(result)

			useCanvas(canvas)

			// create tables
			addTable("track", canvas)
			addTable("user", canvas)
			addTable("label", canvas)

			// return the canvas
			return canvas
		}
	}
}

func useCanvas(canvas *sql.DB) {
	if result, err := canvas.Exec(`USE canvas`); err != nil {
		panic(err)
	} else {
		logResult(result)
	}
}

func AddRow(table string, row map[string]interface{}, canvas *sql.DB) {

	// add columns and values to query
	var columns []string
	var values []interface{}
	for column, value := range row {

		// make sure field isn't empty
		if row[column] != nil && row[column] != "" {

			// check column name
			if column == "user" || column == "label" {
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

				columns = append(columns, entry)
				values = append(values, value)
			}
		}
	}

	query := constructQuery(table, columns)

	// log.Println(query)

	// prepare statment
	if stmt, err := canvas.Prepare(query); err != nil {

		// check for mysql error
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

				addColumn(unknownColumn, table, columnType, canvas)

				// try to add the track again
				AddRow(table, row, canvas)

			} else if prepareErr.Number == 1046 {
				useCanvas(canvas)
			} else {
				panic(err)
			}
		}
	} else {
		defer stmt.Close()

		if result, err := stmt.Exec(values...); err != nil {

			if execErr, ok := err.(*mysql.MySQLError); ok {
				if execErr.Number == 1366 {
					log.Println(values)
					problemColumn := strings.Split(execErr.Message, "'")[3]
					problemValue := strings.Split(execErr.Message, "'")[1]
					log.Println(problemColumn, problemValue, reflect.TypeOf(problemValue))
					row[problemColumn] = problemValue
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {

			logResult(result)
		}
	}
}

func logResult(result sql.Result) {
	if lastID, err := result.LastInsertId(); err != nil {
		panic(err)
	} else {
		if rowsAffected, err := result.RowsAffected(); err != nil {
			panic(err)
		} else {
			log.Printf("Last ID: %d; Rows affected: %d", lastID, rowsAffected)
		}
	}
}

func GetLatest(trackID *int, canvas *sql.DB) {

	if err := canvas.QueryRow(`SELECT MAX(id) FROM track`).
		Scan(trackID); err != nil {
		panic(err)
	}
}
