package util

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysqlResult(username, password, host, project, sql string) (map[string]interface{}, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username,
		password,
		host,
		project,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer d.Close()

	sqlRows, err := db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	columns, err := sqlRows.Columns()
	if err != nil {
		return nil, err
	}
	var rows [][]interface{}
	var cols []*[]interface{}
	for _, _ = range columns {
		var colArray []interface{}
		cols = append(cols, &colArray)
	}
	for sqlRows.Next() {
		m := make(map[string]interface{})
		db.ScanRows(sqlRows, &m)
		var row []interface{}
		for i, field := range columns {
			value := m[field]
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
	sqlRows.Close()
	return resultMap, nil
}
