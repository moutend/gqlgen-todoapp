package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var (
	Db *sql.DB
)

type Option struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func Initialize(opt Option) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		opt.User,
		opt.Password,
		opt.Host,
		opt.Port,
		opt.Database,
	)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	Db = db

	return nil
}

func Migrate(filepath string) error {
	if err := Db.Ping(); err != nil {
		return err
	}

	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		filepath,
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
