package db

import "database/sql"

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

func GetLatest(id *int, table string, canvas *sql.DB) {
	if err := canvas.QueryRow(`SELECT MAX(id) FROM ` + table).
		Scan(id); err != nil {
		if err.Error() == `sql: Scan error on column index 0: converting driver.Value type <nil> ("<nil>") to a int: invalid syntax` {
		} else {
			panic(err)
		}
	}
}
