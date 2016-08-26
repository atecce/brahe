package db

import (
	"database/sql"
	"log"
)

func constructQuery(table string, columns []string) string {

	// insert columns
	query := `INSERT INTO ` + table + `(`
	for _, column := range columns {
		query += column + ", "
	}
	query = query[0 : len(query)-2]

	// insert values
	query += `) VALUES (`
	for _ = range columns {
		query += "?, "
	}
	query = query[0 : len(query)-2]

	// update entire entry on duplicate
	query += `) ON DUPLICATE KEY UPDATE `
	for _, column := range columns {
		query += column + "=VALUES(" + column + "), "
	}
	query = query[0 : len(query)-2]

	return query
}

func logResult(query string, result sql.Result) {
	if lastID, err := result.LastInsertId(); err != nil {
		panic(err)
	} else {
		if rowsAffected, err := result.RowsAffected(); err != nil {
			panic(err)
		} else {
			log.Println(query)
			log.Printf("Last ID: %d; Rows affected: %d", lastID, rowsAffected)
		}
	}
}

func GetLatest(id *int, table string, canvas *sql.DB) {
	if err := canvas.QueryRow(`SELECT MAX(id) FROM ` + table).
		Scan(id); err != nil {
		if err.Error() == `sql: Scan error on column index 0: converting driver.Value type <nil> ("<nil>") to a int: invalid syntax` {
		} else {
			panic(err)
		}
	}
}
