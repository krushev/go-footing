package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
)

type Connection interface {
	GetDB() *gorm.DB
}

type connection struct {
	DB *gorm.DB
}

func (c *connection) GetDB() *gorm.DB {
	return c.DB
}

func NewConnection() Connection {
	user := viper.GetString("db.user")
	pass := viper.GetString("db.pass")
	name := viper.GetString("db.name")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")

	var d gorm.Dialector
	if viper.GetString("db.drv") == "sqlite" {
		d = sqlite.Open(name)
	} else {
		dsn := url.URL{
			User:     url.UserPassword(user, pass),
			Scheme:   "postgres",
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     name,
			RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
		}
		d = postgres.Open(dsn.String())
	}

	db, err := gorm.Open(d, &gorm.Config {
		Logger: logger.Default.LogMode(logMode()),
	})

	if err != nil {
		panic("database connection failed")
	}

	return &connection{DB: db}
}

func logMode() logger.LogLevel {
	dbLogLevel := viper.GetString("db_log")
	if dbLogLevel == "" {
		dbLogLevel = viper.GetString("db.log")
	}
	switch dbLogLevel {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}
	return logger.Silent
}
