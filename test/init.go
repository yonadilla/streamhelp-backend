package test

import (
	"streamhelper-backend/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var App *fiber.App

var DB *gorm.DB

var ViperConfig *viper.Viper

var Log *logrus.Logger

var Validate *validator.Validate

func init(){
	ViperConfig = config.NewViper()
	Log = config.NewLogger(ViperConfig)
	Validate = config.NewValidator(ViperConfig)
	App = config.NewFiber(ViperConfig)
	DB = config.NewDatabase(ViperConfig, Log)
	
	config.Bootstrap(&config.BootstrapConfig{
		DB:       DB,
		App:      App,
		Log:      Log,
		Validate: Validate,
		Config:   ViperConfig,
	})
}