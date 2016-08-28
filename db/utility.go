package db

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func (canvas *Canvas) constructQuery(table string, columns []string) string {

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

func (canvas *Canvas) splitMap(row map[string]interface{}) ([]string, []interface{}) {

	// add columns and values to query
	var columns []string
	var values []interface{}
	for column, value := range row {

		// make sure field isn't empty
		if value != nil && value != "" {

			switch reflect.ValueOf(value).Kind() {

			// create a table for slice TODO
			case reflect.Slice:
				canvas.AddTable(column[:len(column)-1])
				log.Println(column)
				for _, entry := range value.([]interface{}) {
					log.Println(entry)
				}

			// create a new table for additional map
			case reflect.Map:
				canvas.AddTable(column)

				// recursion WTF
				canvas.AddRow(column, value.(map[string]interface{}))
				continue

			default:

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

func (canvas *Canvas) checkMySQLerr(table string, row map[string]interface{}, mysqlErr *mysql.MySQLError) {

	switch mysqlErr.Number {

	// no database selected
	case 1046:
		canvas.use()

	// unknown column
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

		log.Println(columnType)

		// add column and try to add the row again
		canvas.addColumn(unknownColumn, table, columnType)
		canvas.AddRow(table, row)

	// duplicate columns
	case 1060:
		return

	// unknown table
	case 1146:
		canvas.AddTable(table)

	// bananas characters
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

func (canvas *Canvas) GetMissing(table string) map[int]bool {

	// initialize map of missing entries
	missing := make(map[int]bool)

	// for both missing and not missing
	for k, v := range map[string]bool{"": false, "missing_": true} {

		// query the database for the id
		if rows, err := canvas.con.Query(`SELECT id FROM ` + k + table); err != nil {

			// assume error is MySQL specific and try again
			canvas.checkMySQLerr(k+table, nil, err.(*mysql.MySQLError))
			return canvas.GetMissing(table)
		} else {
			defer rows.Close()

			log.Println(table, rows, err)

			// scan ids into map
			for rows.Next() {
				var id int
				rows.Scan(&id)
				missing[id] = v
			}

			err = rows.Err()
			if err != nil {
				panic(err)
			}
		}
	}
	return missing
}

func (canvas *Canvas) AddMissing(method string) {

	// split REST method
	dbInfo := strings.Split(method, "/")

	// table name is first entry without plural
	table := "missing_" + dbInfo[0][0:len(dbInfo[0])-1]

	// row is an id in the second entry
	row := map[string]interface{}{"id": dbInfo[1]}

	// add missing entry
	canvas.AddTable(table)
	canvas.AddRow(table, row)
}
