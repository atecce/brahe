package db

import (
	"database/sql"
	"log"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

type Canvas struct {
	con *sql.DB

	Kind string
	URL  string
	Name string
}

func (canvas *Canvas) use() {

	// use
	query := `USE ` + canvas.Name
	result, err := canvas.con.Exec(query)
	log.Println(query, result)

	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) Initiate() {

	// connect
	con, err := sql.Open(canvas.Kind, canvas.URL)
	canvas.con = con

	// create
	query := `CREATE DATABASE IF NOT EXISTS ` + canvas.Name
	result, err := canvas.con.Exec(query)
	log.Println(query, result)

	// use
	canvas.use()

	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) AddTable(name string) {

	query := `CREATE TABLE IF NOT EXISTS ` + name + ` (id INTEGER NOT NULL, PRIMARY KEY (id))`
	result, err := canvas.con.Exec(query)
	log.Println(query, result)

	if err != nil {
		canvas.checkMySQLerr(name, nil, err.(*mysql.MySQLError))
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
	result, err := canvas.con.Exec(query)
	log.Println(query, result)

	if err != nil {
		canvas.checkMySQLerr(table, nil, err.(*mysql.MySQLError))
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
		result, err := stmt.Exec(values...)
		if err != nil {

			// assert error is MySQL specific
			canvas.checkMySQLerr(table, row, err.(*mysql.MySQLError))

		} else {

			// log only values of insert query
			log.Println("INSERT INTO", columns, "VALUES", values, result)
		}
	}
}
