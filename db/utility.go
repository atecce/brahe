package db

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
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
	query += `) ON DUPLICATE KEY UPDATE `
	for _, column := range columns {
		query += column + "=VALUES(" + column + "), "
	}
	query = query[0 : len(query)-2]

	return query
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

func checkMySQLerr(table string, row map[string]interface{}, canvas *sql.DB, mysqlErr *mysql.MySQLError) {

	switch mysqlErr.Number {

	// handle unknown column
	case 1054:

		// column name is second field delimited with single quotes
		unknownColumn := strings.Split(mysqlErr.Message, "'")[1]

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

	// handle bananas characters
	case 1366:

		// error message convenienty delimted by '
		problemColumn := strings.Split(mysqlErr.Message, "'")[3]
		problemValue := strings.Split(mysqlErr.Message, "'")[1]

		// reset as string
		row[problemColumn] = problemValue

	default:
		panic(mysqlErr)
	}
}

func GetLatest(id *int, table string, canvas *sql.DB) {
	if err := canvas.QueryRow(`SELECT MAX(id) FROM ` + table).
		Scan(id); err != nil {
		if err.Error() == `sql: Scan error on column index 0: converting driver.Value type <nil> ("<nil>") to a int: invalid syntax` {
		} else {
			panic(err)
		}
	}
}
