package db

import (
	"log"
	"reflect"

	"github.com/gocql/gocql"
)

func (canvas *Canvas) constructQuery(table string, columns []string) string {

	// insert columns
	query := `INSERT INTO ` + table + `(`
	for _, column := range columns {
		query += column + ", "
	}
	query = query[:len(query)-2]

	// insert values
	query += `) VALUES (`
	for _ = range columns {
		query += "?, "
	}
	query = query[:len(query)-2] + ")"

	// update entire entry on duplicate
	// query += `) ON DUPLICATE KEY UPDATE `
	// for _, column := range columns {
	// 	query += column + "=VALUES(" + column + "), "
	// }
	// query = query[:len(query)-2]

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
				for _, entry := range value.([]interface{}) {
					log.Println(entry)
				}

			// create a new table for additional map
			// case reflect.Map:
			// 	canvas.AddTable(column)
			//
			// 	// recursion WTF
			// 	canvas.AddRow(column, value.(map[string]interface{}))
			// 	continue

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

func (canvas *Canvas) checkGoCQLerr(table string, row string, err gocql.RequestError) {

	switch err.Code() {

	// // unknown column
	// case 1054:
	//
	// 	// column name is second field delimited with single quotes
	// 	unknownColumn := strings.Split(err.Message(), "'")[1]
	//
	// 	// handle special case for MySQL keyword
	// 	var columnType reflect.Type
	// 	if unknownColumn == "release_number" {
	// 		columnType = reflect.TypeOf(row["release"])
	// 	} else {
	// 		columnType = reflect.TypeOf(row[unknownColumn])
	// 	}
	//
	// 	// add column
	// 	canvas.addColumn(unknownColumn, table, columnType)
	//
	// // duplicate columns
	// case 1060:
	// 	return
	//
	// // bananas characters
	// case 1366:
	//
	// 	// error message convenienty delimted by '
	// 	problemColumn := strings.Split(err.Message(), "'")[3]
	// 	problemValue := strings.Split(err.Message(), "'")[1]
	//
	// 	// reset as string
	// 	row[problemColumn] = problemValue

	// unknown table
	case 8704:
		log.Println(err.Code())
		log.Println(err)
		canvas.AddTable(table)

	default:
		log.Println(err.Code())
		panic(err)
	}
}

// func (canvas *Canvas) GetMissing(table string) map[int]bool {
//
// 	// initialize map of missing entries
// 	missing := make(map[int]bool)
//
// 	// for both missing and not missing
// 	for k, v := range map[string]bool{"": false, "missing_": true} {
//
// 		// query the database for the id
// 		if rows, err := canvas.con.Query(`SELECT id FROM ` + k + table); err != nil {
//
// 			// assume error is MySQL specific and try again
// 			canvas.checkMySQLerr(k+table, nil, err.(*mysql.MySQLError))
// 			return canvas.GetMissing(table)
// 		} else {
// 			defer rows.Close()
//
// 			log.Println(table, rows, err)
//
// 			// scan ids into map
// 			for rows.Next() {
// 				var id int
// 				rows.Scan(&id)
// 				missing[id] = v
// 			}
//
// 			err = rows.Err()
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// 	return missing
// }
