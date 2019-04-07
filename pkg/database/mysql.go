package database

import (
	"encoder-backend/pkg/config"
	"errors"
	"fmt"
	"github.com/Ewan-Walker/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
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

	db, err := gorm.Open("mysql", CreateDSN(username, password, hostname, database, port))
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
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, hostname, port, database)
}

func Migrate(models ...interface{}) {
	if global == nil {
		_, err := Connect()
		if err != nil {
			logrus.WithError(err).Fatal("database: failed to connect")
		}
	}

	global.AutoMigrate(models...)
}
