package db

import (
	"log"
	"reflect"

	"github.com/gocql/gocql"
)

type Canvas struct {
	Kind    string
	IP      string
	Name    string
	Session *gocql.Session
}

func (canvas *Canvas) Initiate() {

	cluster := gocql.NewCluster(canvas.IP)
	cluster.ProtoVersion = 4
	cluster.Keyspace = canvas.Name
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	canvas.Session = session
}

func (canvas *Canvas) AddTable(name string) {

	query := `CREATE TABLE IF NOT EXISTS ` + name + ` (id INT, PRIMARY KEY (id))`
	err := canvas.Session.Query(query).Exec()
	if err != nil {
		panic(err)
	}
	log.Println(query)
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
	err := canvas.Session.Query(query).Exec()
	if err != nil {
		panic(err)
	}
	log.Println(query)
}

func (canvas *Canvas) AddRow(table string, row string) {

	// // split map into columns and values
	// columns, values := canvas.splitMap(row)
	//
	// // construct query out of lists
	// query := canvas.constructQuery(table, columns)

	query := "INSERT INTO " + table + " JSON `" + row + "`"

	log.Println(query)

	// prepare statment
	stmt := canvas.Session.Query(query)

	log.Println(stmt)

	// insert row
	err := stmt.Exec()
	if err != nil {

		// assert error is MySQL specific
		canvas.checkGoCQLerr(table, row, err.(gocql.RequestError))
		canvas.AddRow(table, row)

	} else {

		// log only values of insert query
		log.Println(query)
		// log.Println("INSERT INTO", table, columns, "VALUES", values)
	}
}

// func (canvas *Canvas) AddMissing(method string) {
//
// 	// split REST method
// 	dbInfo := strings.Split(method, "/")
//
// 	// table name is first entry without plural
// 	table := "missing_" + dbInfo[0][:len(dbInfo[0])-1]
//
// 	// row is an id in the second entry
// 	row := map[string]interface{}{"id": dbInfo[1]}
//
// 	// add missing entry
// 	canvas.AddTable(table)
// 	canvas.AddRow(table, row)
// }
