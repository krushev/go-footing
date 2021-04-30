package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/krushev/go-footing/db"
	"github.com/krushev/go-footing/models"
	"github.com/krushev/go-footing/security"
	"github.com/krushev/go-footing/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewApp(debug bool) *fiber.App {
	conn := db.NewConnection()
	// Auto Migrate
	conn.GetDB().AutoMigrate(&models.User{})

	us := services.NewUserService(conn)
	bs := services.NewBootstrapService(us)
	bs.InitUsers()

	uc := NewUserController(us)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api")
	api001 := api.Group("/v0.0.1")

	jwtBus := security.InitAuth(us)

	api.Post("/login", jwtBus.LoginHandler)
	api.Post("/logout", jwtBus.LogoutHandler)
	api.Post("/refresh", jwtBus.RefreshHandler)

	if !debug {
		api001.Use(jwtBus.MiddlewareFunc())
	}

	api001.Mount("/users", uc.Router())

	port := viper.GetString("app.port")
	zap.S().Infof("app.port: %s", port)

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("pong"))
	})

	if !debug {
		zap.S().Fatal(app.Listen(":" + port))
	}

	return app
}
