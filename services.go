package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func PingDB(url string) (*gorm.DB, bool) {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	DB = db
	return DB, err == nil
}

func GetTables() ([]string, bool) {
	var tables []string

	DB.Raw("SHOW TABLES").Scan(&tables)

	if len(tables) == 0 {
		return nil, false
	}

	return tables, true
}

func GetColumns(table string) []string {

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
