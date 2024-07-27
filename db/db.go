package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetCurrentTime() (string,error) {
	var currentTime string

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return "sql.Open error", err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	
	err = db.QueryRow(`SELECT datetime('now', '+9 hours')`).Scan(&currentTime)
	if err != nil {
		return "db.QueryRow error", err
	} 
	return currentTime, nil
}

func DBConn() (*bun.DB,error) {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	// Open a PostgreSQL database.
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",dbUser,dbPass,dbHost,dbPort,dbName)
	fmt.Println(dsn)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it.
	db := bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db, nil

}