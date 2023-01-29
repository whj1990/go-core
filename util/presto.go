package util

import (
	dbsql "database/sql"
	"fmt"
	"strings"
)

func GetPrestoResult(username, password, host, project, sql string) (map[string]interface{}, error) {
	dsn := fmt.Sprintf("http://%s:%s@%s&schema=%s",
		username,
		password,
		host,
		project)
	db, _ := dbsql.Open("presto", dsn)
	resultRows, _ := db.Query(strings.ReplaceAll(sql, ";", ""))
	columns, _ := resultRows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	//var rows []map[string]interface{} //返回的切片
	var rows []interface{}
	var cols []*[]interface{}
	for _, _ = range columns {
		var colArray []interface{}
		cols = append(cols, &colArray)
	}
	for resultRows.Next() {
		_ = resultRows.Scan(cache...)
		//item := make(map[string]interface{})
		var row []interface{}
		for i, data := range cache {
			value := *data.(*interface{})
			//item[columns[i]] = *data.(*interface{}) //取实际类型
			row = append(row, value)
			col := cols[i]
			*col = append(*col, value)
		}
		rows = append(rows, row)
	}
	resultMap := make(map[string]interface{})
	resultMap["fields"] = columns
	resultMap["rows"] = rows
	resultMap["cols"] = cols
	resultRows.Close()
	return resultMap, nil
}
