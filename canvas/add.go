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
	Client *bigtable.Client
	AC     *bigtable.AdminClient
}

func (canvas *Canvas) Initiate() {

	// create admin client for adding tables and families
	ac, err := bigtable.NewAdminClient(context.Background(), "telos-143019", "uraniborg")
	if err != nil {
		panic(err)
	}
	canvas.AC = ac

	// create normal client for adding entries
	client, err := bigtable.NewClient(context.Background(), "telos-143019", "uraniborg")
	if err != nil {
		panic(err)
	}
	canvas.Client = client
}

func (canvas *Canvas) AddTable(table string) {
	err := canvas.AC.CreateTable(context.Background(), table)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}
}

func (canvas *Canvas) AddFamily(table, family string) {
	err := canvas.AC.CreateColumnFamily(context.Background(), table, family)
	switch grpc.Code(err) {
	case codes.OK, codes.AlreadyExists:
	default:
		panic(err)
	}
}

func (canvas *Canvas) Record(table, row, family, column string) {

	// open table
	tbl := canvas.Client.Open(table)

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
	default:
		panic(err)
	}
}
