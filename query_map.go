package gotools

// Query 用于快速执行 SQL 语句，返回 []string（Columns 列名）、[][]interface{}（各行对应的值），对 value 操作时，需要知道其字段对应的类型。
// Usage:
// _, a, _ := Query(dbPtr, "select first_name from customers where balance >= ?", 1000)
//
//	for i := range a {
//		fmt.Println(a[i].(string)) // you must know what type the db driver converts your columns to
//	}
func Query(db DB, sql string, args ...interface{}) ([]string, [][]interface{}, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, nil, err
	}

	// use rows.Columns() to get a reference to all column names in the result
	colNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()
	res := make([][]interface{}, 0)
	for rows.Next() {
		columnPointers := make([]interface{}, len(colNames))
		// This is necessary because the sql package requires pointers when scanning
		for i := range columnPointers {
			columnPointers[i] = new(interface{})
		}

		// Scan the result into the column pointers...
		err = rows.Scan(columnPointers...)
		if err != nil {
			return nil, nil, err
		}

		values := make([]interface{}, len(colNames))
		for idx := range colNames {
			values[idx] = *columnPointers[idx].(*interface{})
		}
		res = append(res, values)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return colNames, res, nil
}

// QueryMap 用于快速执行 SQL 语句，返回 []string（Columns 列名）、[]map[string]interface{}（各列对应的值），对 value 操作时，需要知道其字段对应的类型。
// Usage:
// _, a, _ := QueryMap(dbPtr, "select first_name from customers where balance >= ?", 1000)
//
//	for _, firstName := range a {
//		fmt.Println(a["first_name"].(string)) // you must know what type the db driver converts your columns to
//	}
func QueryMap(db DB, sql string, args ...interface{}) ([]string, []map[string]interface{}, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, nil, err
	}

	// use rows.Columns() to get a reference to all column names in the result
	colNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()
	res := make([]map[string]interface{}, 0)
	for rows.Next() {
		columnPointers := make([]interface{}, len(colNames))
		// This is necessary because the sql package requires pointers when scanning
		for i := range columnPointers {
			columnPointers[i] = new(interface{})
		}

		// Scan the result into the column pointers...
		err = rows.Scan(columnPointers...)
		if err != nil {
			return nil, nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		these := make(map[string]interface{})
		for idx, colName := range colNames {
			these[colName] = *columnPointers[idx].(*interface{})
		}
		res = append(res, these)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return colNames, res, nil
}

// QueryMapOne 是 QueryMap 的单个结果。
func QueryMapOne(db DB, sql string, args ...interface{}) ([]string, map[string]interface{}, error) {
	columns, xs, err := QueryMap(db, sql, args...)
	if len(xs) == 0 {
		return columns, nil, err
	}
	return columns, xs[0], nil
}
