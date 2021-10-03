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

func execMigrations() error {
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
func InitDB(iLog *log.Logger) error {
	// creating db file or getting access to it
	iLog.Println("accessing database file")
	ConnToDB, _ = sql.Open("sqlite3", "file:db/alber.db?_auth&_auth_user=alber&_auth_pass=zhibek1234&_auth_crypt=sha1")
	iLog.Println("access completed")

	// make some sql settings
	iLog.Println("set up database configs")
	if _, e := ConnToDB.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;PRAGMA case_sensitive_like = true;PRAGMA auto_vacuum = FULL;"); e != nil {
		return e
	}
	iLog.Println("database configured")

	// do migrations
	iLog.Println("making migrations")
	if e := execMigrations(); e != nil {
		return e
	}
	iLog.Println("migrations completed")
	return nil
}
