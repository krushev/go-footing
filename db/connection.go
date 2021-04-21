package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
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

	var enableLogging = logger.Default.LogMode(logger.Info)
	if !viper.GetBool("db_debug") && !viper.GetBool("db.debug") {
		enableLogging.LogMode(logger.Silent)
	}

	dsn := url.URL{
		User:     url.UserPassword(user, pass),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     name,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config {
		Logger: enableLogging,
	})

	if err != nil {
		panic("database connection failed")
	}

	return &connection{DB: db}
}
