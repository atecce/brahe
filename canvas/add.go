package canvas

import (
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
	client *bigtable.Client
	ac     *bigtable.AdminClient
}

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

func (canvas *Canvas) Close() {
	canvas.client.Close()
	canvas.ac.Close()
}

func (canvas *Canvas) AddTable(table string) {
	err := canvas.ac.CreateTable(context.Background(), table)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}
}

func (canvas *Canvas) AddFamily(table, family string) {
	err := canvas.ac.CreateColumnFamily(context.Background(), table, family)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}
}

func (canvas *Canvas) Record(table, row, family, column string) {

	// open table
	tbl := canvas.client.Open(table)

	// set mutation
	mut := bigtable.NewMutation()
	mut.Set(family, column, bigtable.ServerTime, []byte("1"))

	// add entry
	err := tbl.Apply(context.Background(), row, mut)
	switch grpc.Code(err) {
	case codes.OK:
		return
	case codes.NotFound:
		canvas.AddFamily(table, family)
	case codes.Internal:
		canvas.Initiate()
	default:
		panic(err)
	}
}
