package orm

import (
	"context"
	"database/sql"
	"log"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"

	// db migration
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
)

// ConnToDB is conn to db
var ConnToDB *sql.DB

func checkFatal(eLogger *log.Logger, e error) {
	if e != nil {
		eLogger.Fatal(e)
	}
}

func execMigrations(eLogger *log.Logger) error {
	driver, e := sqlite3.WithInstance(ConnToDB, &sqlite3.Config{})
	if e != nil {
		return e
	}

	m, e := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
	if e != nil {
		return e
	}

	if e = m.Up(); e != nil && e != migrate.ErrNoChange {
		return e
	}
	return nil
}

// InitDB init db, settings and tables
func InitDB(eLog, iLog *log.Logger) {
	// creating db file or getting access to it
	iLog.Println("accessing database file")
	ConnToDB, _ = sql.Open("sqlite3", "file:db/zhibek.db?_auth&_auth_user=zhibek&_auth_pass=zhibek1234&_auth_crypt=sha1")
	iLog.Println("access completed")

	// make some sql settings
	iLog.Println("set up database configs")
	_, e := ConnToDB.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;PRAGMA case_sensitive_like = true;PRAGMA auto_vacuum = FULL;")
	checkFatal(eLog, e)
	iLog.Println("database configured")

	// do migrations
	iLog.Println("making migrations")
	checkFatal(eLog, execMigrations(eLog))
	iLog.Println("migrations completed")
}
