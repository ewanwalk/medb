package database

import (
	"errors"
	"fmt"
	"github.com/Ewan-Walker/gorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	global                *gorm.DB
	ErrDatabaseConnection = errors.New("db: connection closed")
)

func Connect() (*gorm.DB, error) {

	if global != nil {
		return global, nil
	}

	// TODO make dsn from env
	dsn := "root:master1@tcp(localhost:3306)/medb_dev_2?parseTime=true"

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxOpenConns(25)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

	global = db

	return global, nil
}

func CreateDSN(username, password, hostname, database, port string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, hostname, database, port)
}

func Migrate(models ...interface{}) {
	if global == nil {
		Connect()
	}

	global.AutoMigrate(models...)
}
