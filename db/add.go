package db

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

const (
	project  = "telos-143019"
	instance = "uraniborg"
)

type Canvas struct {
	Kind   string
	IP     string
	Name   string
	client *bigtable.Client
	ac     *bigtable.AdminClient
}

// TODO maybe close these
func (canvas *Canvas) Initiate() {

	// create admin client for adding tables and families
	ac, err := bigtable.NewAdminClient(context.Background(), project, instance)
	if err != nil {
		panic(err)
	}
	canvas.ac = ac

	// create normal client for adding entries
	client, err := bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		panic(err)
	}
	canvas.client = client
}

func (canvas *Canvas) AddTable(name string) {

	err := canvas.ac.CreateTable(context.Background(), name)
	if err != nil {

		if grpc.Code(err) == codes.AlreadyExists {
			return
		}

		panic(err)
	}
}

func (canvas *Canvas) addFamily(table, family string) {
	err := canvas.ac.CreateColumnFamily(context.Background(), table, family)
	if err != nil {

		switch grpc.Code(err) {

		case codes.NotFound:
			canvas.AddTable(table)
			canvas.addFamily(table, family)

		case codes.AlreadyExists:
			return

		default:
			panic(err)
		}
	}
}

func (canvas *Canvas) AddRow(table string, row map[string]interface{}) {

	// get id into big endian
	id := []byte(strconv.FormatFloat(row["id"].(float64), 'f', -1, 64))

	log.Println("id: ", row["id"], id)

	// TODO ApplyBulk (looks like not needed)
	for family, column := range row {

		// skip id family and nil column
		if family == "id" || column == nil {
			continue
		}

		log.Println("table: ", family)
		log.Println("family: ", family, reflect.ValueOf(family).Kind())
		log.Println("column: ", column, reflect.ValueOf(column).Kind())
		log.Println()

		var typedColumn string
		switch reflect.ValueOf(column).Kind() {

		// TODO handle subscriptions
		case reflect.Slice:
			continue

		// recursively walk down nested tree
		case reflect.Map:
			canvas.AddRow(family, column.(map[string]interface{}))

		case reflect.Float64:
			typedColumn = strconv.FormatFloat(column.(float64), 'f', -1, 64)

		case reflect.Bool:
			typedColumn = strconv.FormatBool(column.(bool))

		default:
			typedColumn = column.(string)
		}
		canvas.addColumn(id, family, typedColumn, table)
	}
}

func (canvas *Canvas) addColumn(id []byte, family, typedColumn, table string) {

	mut := bigtable.NewMutation()
	mut.Set(family, typedColumn, bigtable.ServerTime, id)

	// open table
	tbl := canvas.client.Open(table)

	// add column
	err := tbl.Apply(context.Background(), string(id), mut)
	if err != nil {

		switch grpc.Code(err) {

		case codes.InvalidArgument:
			panic(err)

		// add family and try again if family not found
		case codes.NotFound:
			canvas.addFamily(table, family)
			canvas.addColumn(id, family, typedColumn, table)
			return
		}

		panic(err)
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
