package db

// func (canvas *Canvas) checkGoCQLerr(table string, row map[string]interface{}, err gocql.RequestError) {
//
// 	switch err.Code() {
//
// 	// errInvalid in gocql
// 	case 8704:
//
// 		log.Println(err.Error())
//
// 		// regexes to parse error message
// 		missingTable := regexp.MustCompile(`^unconfigured table (.*)$`)
// 		missingColumn := regexp.MustCompile(`^JSON values map contains unrecognized column: (.*)$`)
//
// 		// check for missing table
// 		if missingTable.MatchString(err.Error()) {
//
// 			// get missing table name from error message
// 			missing := missingTable.FindStringSubmatch(err.Error())[1]
//
// 			// add the table
// 			canvas.AddTable(missing)
//
// 			// check for missing column
// 		} else if missingColumn.MatchString(err.Error()) {
//
// 			// get missing column from error message
// 			column := missingColumn.FindStringSubmatch(err.Error())[1]
//
// 			// get type info
// 			columnType := reflect.TypeOf(row[column])
// 			log.Println(column)
// 			log.Println(columnType)
// 			if columnType != nil {
// 				canvas.addColumn(column, table, columnType)
// 			}
//
// 		} else {
//
// 			log.Println(row)
// 			log.Println(err.Code())
// 			panic(err)
// 		}
//
// 	default:
// 		log.Println(err.Code())
// 		panic(err)
// 	}
// }
//
// func (canvas *Canvas) GetMissing(table string) map[int]bool {
//
// 	// initialize map of missing entries
// 	missing := make(map[int]bool)
//
// 	// for both missing and not missing
// 	for k, v := range map[string]bool{"": false, "missing_": true} {
//
// 		// query the database for the id
// 		rows := canvas.Session.Query(`SELECT id FROM ` + k + table).Iter()
//
// 		log.Println(`SELECT id FROM ` + k + table)
//
// 		// scan ids into map
// 		var id int
// 		for rows.Scan(&id) {
// 			missing[id] = v
// 			log.Println(id, v)
// 		}
//
// 		if err := rows.Close(); err != nil {
//
// 			// assume error is MySQL specific and try again
// 			canvas.checkGoCQLerr(k+table, nil, err.(gocql.RequestError))
// 			return canvas.GetMissing(table)
// 		}
// 	}
// 	return missing
// }
