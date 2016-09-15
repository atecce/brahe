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
		panic(err)
	}
}

func (canvas *Canvas) AddRow(table string, row map[string]interface{}) {

	// get id into big endian
	id := []byte(strconv.FormatFloat(row["id"].(float64), 'f', -1, 64))

	log.Println("id: ", row["id"], id)

	// TODO reflection and ApplyBulk
	for family, column := range row {

		// skip id family and nil column
		if family == "id" || column == nil {
			continue
		}

		log.Println("family: ", family, reflect.ValueOf(family).Kind())
		log.Println("column: ", column, reflect.ValueOf(column).Kind())
		log.Println()

		mut := bigtable.NewMutation()

		switch reflect.ValueOf(column).Kind() {

		// TODO handle subscriptions
		case reflect.Slice:
			continue

		// recursively walk down nested tree
		case reflect.Map:
			canvas.AddRow(table, row)

		case reflect.Float64:
			typedColumn := strconv.FormatFloat(column.(float64), 'f', -1, 64)
			mut.Set(family, string(typedColumn), bigtable.ServerTime, id)

		case reflect.Bool:
			typedColumn := strconv.FormatBool(column.(bool))
			mut.Set(family, string(typedColumn), bigtable.ServerTime, id)

		default:

			mut.Set(family, column.(string), bigtable.ServerTime, id)

		}
	}
}

func (canvas *Canvas) addColumn(id []byte, family, table string, mut *bigtable.Mutation) {

	// open table
	tbl := canvas.client.Open(table)

	// add column
	err := tbl.Apply(context.Background(), string(id), mut)
	if err != nil {

		// add family and try again if family not found
		if grpc.Code(err) == codes.NotFound {
			canvas.addFamily(table, family)
			canvas.addColumn(id, family, table, mut)
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
