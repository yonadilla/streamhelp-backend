package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	dbname := viper.GetString("database.dbname")
    user := viper.GetString("database.user")
    password := viper.GetString("database.password")
    sslmode := viper.GetString("database.sslmode")
    timezone := viper.GetString("database.timezone")
    host := fmt.Sprintf("%s", viper.Get("database.host")) 
    port := fmt.Sprintf("%v", viper.Get("database.port"))
    timeout := fmt.Sprintf("%v", viper.Get("database.connect_timeout"))
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s connect_timeout=%s TimeZone=%s",
        host, user, password, dbname, port, sslmode, timeout, timezone)
	db , err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold: time.Second * 5,
			Colorful: false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries: true,
			LogLevel: logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect database : %v ", err)
	}

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}