package config

import (
	"streamhelper-backend/internal/delivery/http"
	"streamhelper-backend/internal/delivery/http/middleware"
	"streamhelper-backend/internal/delivery/http/route"
	"streamhelper-backend/internal/entity"
	"streamhelper-backend/internal/repository"
	"streamhelper-backend/internal/usecase"
	"streamhelper-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB 			*gorm.DB
	App 		*fiber.App
	Log			*logrus.Logger
	Validate	*validator.Validate
	Config 		*viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	//setup repository
	userRepository := repository.NewUserRepository(config.Log)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
	})

	tokenUtil := util.NewTokenUtil("benar, benar, rahasia", redisClient)

	// setup use cases
	userUseCase := usecase.NewUserUserCase(config.DB, config.Log, config.Validate, userRepository, tokenUtil)
	
	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	err := config.DB.AutoMigrate(&entity.User{})
    if err != nil {
        config.Log.Fatalf("Gagal migrasi database: %v", err)
    }
	routeConfig := route.RouteConfig{
		App: config.App,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()
}