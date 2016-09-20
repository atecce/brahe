package db

import (
	"log"
	"reflect"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

const (
	project     = "telos-143019"
	instance    = "uraniborg"
	bigtableMax = 16384
)

type Canvas struct {
	Name   string
	client *bigtable.Client
	ac     *bigtable.AdminClient
}

func (canvas *Canvas) Initiate() {

	// create admin client for adding tables and families
	// TODO maybe close these
	ac, err := bigtable.NewAdminClient(context.Background(), project, instance)
	if err != nil {
		panic(err)
	}
	canvas.ac = ac

	// create normal client for adding entries
	// TODO maybe close these
	client, err := bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		panic(err)
	}
	canvas.client = client

	// initialize table
	err = canvas.ac.CreateTable(context.Background(), canvas.Name)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}

	// initialize families
	canvas.AddFamily("follows")
	canvas.AddFamily("favorites")
}

func (canvas *Canvas) AddFamily(family string) {
	err := canvas.ac.CreateColumnFamily(context.Background(), canvas.Name, family)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}
}

func (canvas *Canvas) AddEntry(row, family string, entry map[string]interface{}) {

	// TODO ApplyBulk (looks like not needed)
	for column, typedValue := range entry {

		// skip and nil values
		if typedValue == nil {
			continue
		}

		// TODO need more intelligent logging
		log.Println("family: ", family)
		log.Println("column: ", column, reflect.ValueOf(column).Kind())
		log.Println("value: ", typedValue, reflect.ValueOf(typedValue).Kind())
		log.Println()

		// convert value to string
		var value []byte
		switch typedValue.(type) {

		// recursively walk down nested tree
		case map[string]interface{}:
			canvas.AddEntry(row, column, typedValue.(map[string]interface{}))

		// handle basic data types
		case float64:
			value = []byte(strconv.FormatFloat(typedValue.(float64), 'f', -1, 64))
		case bool:
			value = []byte(strconv.FormatBool(typedValue.(bool)))

		// make sure string is below bigtable's maximum length
		case string:
			value = []byte(typedValue.(string))
			if len(value) > bigtableMax {
				// TODO need more intelligent logging
				log.Printf("Value %s in column %s too long", value, column)
				value = value[:bigtableMax]
			}

		default:
			panic(typedValue)
		}

		// add column
		// TODO need more intelligent logging
		log.Println("Adding entry")
		canvas.addColumn(row, family, column, value)
	}
}

func (canvas *Canvas) addColumn(row, family, column string, value []byte) {

	// open table
	tbl := canvas.client.Open(canvas.Name)

	// set mutation
	mut := bigtable.NewMutation()
	mut.Set(family, column, bigtable.ServerTime, value)

	// add column
	err := tbl.Apply(context.Background(), row, mut)
	switch grpc.Code(err) {
	case codes.OK:
		return
	case codes.NotFound:
		canvas.AddFamily(family)
	default:
		panic(err)
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
// 	// entry is an id in the second entry
// 	entry := map[string]interface{}{"id": dbInfo[1]}
//
// 	// add missing entry
// 	canvas.AddTable(table)
// 	canvas.AddEntry(table, entry)
// }
