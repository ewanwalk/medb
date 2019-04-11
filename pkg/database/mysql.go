package database

import (
	"encoder-backend/pkg/config"
	"errors"
	"fmt"
	"github.com/Ewan-Walker/gorm"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
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

	username := ""
	if env := os.Getenv(config.EnvDBUsername); len(env) != 0 {
		username = env
	}

	password := ""
	if env := os.Getenv(config.EnvDBPassword); len(env) != 0 {
		password = env
	}

	hostname := ""
	if env := os.Getenv(config.EnvDBHostname); len(env) != 0 {
		hostname = env
	}

	port := ""
	if env := os.Getenv(config.EnvDBPort); len(env) != 0 {
		port = env
	}

	database := ""
	if env := os.Getenv(config.EnvDBName); len(env) != 0 {
		database = env
	}

	db, err := gorm.Open("mysql", CreateDSN(username, password, hostname, "", port))
	if err != nil {
		return nil, err
	}

	// This could allow for SQL injection should someone have local access
	qry := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", database)

	// create a db if it doesnt exist
	if err := db.Exec(qry).Error; err != nil {
		return nil, err
	}

	db.Close()

	db, err = gorm.Open("mysql", CreateDSN(username, password, hostname, database, port))
	if err != nil {
		return nil, err
	}

	// TODO determine if we prefer this
	db.SetLogger(log.StandardLogger())

	db.DB().SetMaxOpenConns(25)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

	if env := os.Getenv(config.EnvDBDebug); len(env) != 0 {
		db = db.Debug()
	}

	global = db

	return global, nil
}

func CreateDSN(username, password, hostname, database, port string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", username, password, hostname, port, database)
}

func Migrate(models ...interface{}) {
	if global == nil {
		_, err := Connect()
		if err != nil {
			log.WithError(err).Fatal("database: failed to connect")
		}
	}

	global.AutoMigrate(models...)
}
