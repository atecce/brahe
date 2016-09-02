package db

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/gocql/gocql"
)

type Canvas struct {
	Kind    string
	IP      string
	Name    string
	Session *gocql.Session
}

func (canvas *Canvas) Initiate() {

	// create cluster
	cluster := gocql.NewCluster(canvas.IP)
	cluster.ProtoVersion = 4
	cluster.Keyspace = canvas.Name

	// start session
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
	goToCQL := map[string]string{
		"bool":    "BOOLEAN",
		"float64": "FLOAT",
		"string":  "TEXT",
	}

	// add column name and type
	query := `ALTER TABLE ` + table + ` ADD ` + column + ` ` + goToCQL[columnType.String()]
	err := canvas.Session.Query(query).Exec()
	if err != nil {
		canvas.checkGoCQLerr(table, nil, err.(gocql.RequestError))
	}
	log.Println(query)
}

func (canvas *Canvas) AddRow(table string, row map[string]interface{}) {

	for k, v := range row {

		if v != nil {

			switch reflect.ValueOf(v).Kind() {

			// create a table for slice TODO
			case reflect.Slice:
				delete(row, k)
				for _, entry := range v.([]interface{}) {
					log.Println(entry)
				}

			// create a new table for additional map
			case reflect.Map:
				canvas.AddTable(k)
				delete(row, k)

				// recursion WTF
				canvas.AddRow(k, v.(map[string]interface{}))
				continue
			}
		} else {
			delete(row, k)
		}
	}

	raw, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}

	// prepare statment
	query := `INSERT INTO ` + table + ` JSON ?`
	stmt := canvas.Session.Query(query, raw)

	// insert row
	err = stmt.Exec()
	if err != nil {

		// fix error and try again
		canvas.checkGoCQLerr(table, row, err.(gocql.RequestError))
		canvas.AddRow(table, row)

	} else {

		// log only values of insert query
		log.Println(query)
		// log.Println("INSERT INTO", table, columns, "VALUES", values)
	}
}

func (canvas *Canvas) AddMissing(method string) {

	// split REST method
	dbInfo := strings.Split(method, "/")

	// table name is first entry without plural
	table := "missing_" + dbInfo[0][:len(dbInfo[0])-1]

	// row is an id in the second entry
	row := map[string]interface{}{"id": dbInfo[1]}

	// add missing entry
	canvas.AddTable(table)
	canvas.AddRow(table, row)
}
