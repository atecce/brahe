package db

import (
	"bodhi/herodotus"
	"database/sql"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

var dbLog = herodotus.CreateFileLog("db")

type Canvas struct {
	con *sql.DB

	Kind string
	URL  string
	Name string
}

func (canvas *Canvas) Initiate() {

	// connect
	con, err := sql.Open(canvas.Kind, canvas.URL)
	canvas.con = con

	// create
	query := `CREATE DATABASE IF NOT EXISTS ` + canvas.Name
	_, err = canvas.con.Exec(query)
	dbLog.Println(query)

	// use
	query = `USE ` + canvas.Name
	_, err = canvas.con.Exec(query)
	dbLog.Println(query)

	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) AddTable(name string) {

	query := `CREATE TABLE IF NOT EXISTS ` + name + ` (id INTEGER NOT NULL, PRIMARY KEY (id))`
	result, err := canvas.con.Exec(query)
	dbLog.Println(query, result)

	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) addColumn(column, table string, columnType reflect.Type) {

	// map for conversion between go and mysql data types
	goToMySQL := map[string]string{
		"bool":    "BOOL",
		"float64": "FLOAT",
		"string":  "TEXT",
	}

	// add column name and type
	query := `ALTER TABLE ` + table + ` ADD ` + column + ` ` + goToMySQL[columnType.String()]
	_, err := canvas.con.Exec(query)
	dbLog.Println(query)

	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) AddRow(table string, row map[string]interface{}) {

	// split map into columns and values
	columns, values := canvas.splitMap(row)

	// construct query out of lists
	query := canvas.constructQuery(table, columns)

	// prepare statment
	stmt, err := canvas.con.Prepare(query)
	if err != nil {

		// assert error is MySQL specific
		canvas.checkMySQLerr(table, row, err.(*mysql.MySQLError))

	} else {
		defer stmt.Close()

		// insert row
		_, err := stmt.Exec(values...)
		if err != nil {

			// assert error is MySQL specific
			canvas.checkMySQLerr(table, row, err.(*mysql.MySQLError))

		} else {

			// log only values of insert query
			dbLog.Println("INSERT INTO", columns, "VALUES", values)
		}
	}
}
