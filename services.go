package main

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type App struct{}

func (a App) PingDB(url string) (*gorm.DB, bool) {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	DB = db
	return DB, err == nil
}

func (a App) GetTables() ([]string, bool) {
	var tables []string

	DB.Raw("SHOW TABLES").Scan(&tables)

	if len(tables) == 0 {
		return nil, false
	}

	return tables, true
}

func (a App) GetColumns(table string) []string {

	query := fmt.Sprintf("SELECT * FROM %s", table)

	rows, err := DB.Raw(query).Rows()
	if err != nil {
		return nil
	}

	defer rows.Close()

	// scan column
	columns, err := rows.Columns()
	if err != nil {
		return nil
	}

	return columns
}

func (a App) GetValues(table string) []map[string]interface{} {
	var result []map[string]interface{}
	columns := a.GetColumns(table)

	joinColumns := strings.Join(columns, ",")

	query := fmt.Sprintf("SELECT %s FROM %s", joinColumns, table)

	rows, err := DB.Raw(query).Rows()

	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		rowData := make(map[string]interface{})

		values := make([]interface{}, len(columns))
		valuesPtrs := make([]interface{}, len(columns))

		for i := range valuesPtrs {
			valuesPtrs[i] = &values[i]
		}

		if err := rows.Scan(valuesPtrs...); err != nil {
			return nil
		}

		for i, col := range columns {
			val := values[i]
			if val != nil {
				if databytes, ok := val.([]uint8); ok {
					rowData[col] = string(databytes)
				} else {
					rowData[col] = val
				}
			} else {
				rowData[col] = nil
			}
		}

		result = append(result, rowData)

	}

	return result
}
