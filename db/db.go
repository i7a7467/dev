package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func GetCurrentTime() (string,error) {
	var currentTime string

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return "sql.Open error", err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	
	err = db.QueryRow(`SELECT datetime('now')`).Scan(&currentTime)
	if err != nil {
		return "db.QueryRow error", err
	} 
	return currentTime, nil
}