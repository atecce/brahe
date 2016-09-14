package db

import (
	"reflect"
	"strings"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

type Canvas struct {
	Kind   string
	IP     string
	Name   string
	client *bigtable.Client
	ac     *bigtable.AdminClient
}

func (canvas *Canvas) Initiate() {

	// create admin client for adding tables and families
	ac, err := bigtable.NewAdminClient(context.Background(), "telos-143019", "uraniborg")
	if err != nil {
		panic(err)
	}
	canvas.ac = ac

	// create normal client for adding entries
	client, err := bigtable.NewClient(context.Background(), "telos-143019", "uraniborg")
	if err != nil {
		panic(err)
	}
	canvas.client = client
}

func (canvas *Canvas) AddTable(name string) {

	err := canvas.ac.CreateTable(context.Background(), name)
	if err != nil {
		panic(err)
	}

}

func (canvas *Canvas) addFamily(table, family string, columnType reflect.Type) {
	err := canvas.ac.CreateColumnFamily(context.Background(), table, family)
	if err != nil {
		panic(err)
	}
}

func (canvas *Canvas) AddRow(table string, row map[string]interface{}) {

	// set id
	id := row["id"]

	// TODO reflection and ApplyBulk
	for family, column := range row {

		if family == "id" {
			continue
		}

		mut := bigtable.NewMutation()
		mut.Set(family, column.(string), bigtable.ServerTime, id.([]byte))

		tbl := canvas.client.Open(table)
		err := tbl.Apply(context.Background(), id.(string), mut)
		if err != nil {
			panic(err)
		}

	}

	// for k, v := range row {
	//
	// 	if v != nil {
	//
	// 		switch reflect.ValueOf(v).Kind() {
	//
	// 		// create a table for slice TODO
	// 		case reflect.Slice:
	// 			delete(row, k)
	// 			for _, entry := range v.([]interface{}) {
	// 				log.Println(entry)
	// 			}
	//
	// 		// create a new table for additional map
	// 		case reflect.Map:
	// 			canvas.AddTable(k)
	// 			delete(row, k)
	//
	// 			// recursion WTF
	// 			canvas.AddRow(k, v.(map[string]interface{}))
	// 			continue
	// 		}
	// 	} else {
	// 		delete(row, k)
	// 	}
	// }
	//
	// raw, err := json.Marshal(row)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // prepare statment
	// query := `INSERT INTO ` + table + ` JSON ?`
	// stmt := canvas.Session.Query(query, raw)
	//
	// // insert row
	// err = stmt.Exec()
	// if err != nil {
	//
	// 	// fix error and try again
	// 	canvas.checkGoCQLerr(table, row, err.(gocql.RequestError))
	// 	canvas.AddRow(table, row)
	//
	// } else {
	//
	// 	// log only values of insert query
	// 	log.Println(query)
	// 	// log.Println("INSERT INTO", table, columns, "VALUES", values)
	// }
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
