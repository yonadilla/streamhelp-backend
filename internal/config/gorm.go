package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	port := viper.GetString("database.port")
	sslmode := viper.GetString("database.sslmode")
	connect_timeout := viper.GetString("database.connect_timeout")
	timezone := viper.GetString("database.timezone")

	dsn := fmt.Sprintf( "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s connect_timeout=%s TimeZone=%s",
	host,
	user,
	password,
	dbname,
	port,
	sslmode,
	connect_timeout,
	timezone,
	)

	db , err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database : %v ", err)
	}

	return  db
}