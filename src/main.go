package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"api_fiber/src/database"
	handlers_v1 "api_fiber/src/handlers/v1"
	"api_fiber/src/middlewares"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err := database.InitDatabase(
		os.Getenv("HOST_DB"),
		os.Getenv("PORT_DB"),
		os.Getenv("NAME_DB"),
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("SSLMODE_DB"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	api := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Learn fiber v1.0.0",
	})

	api.Static("/", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	api.Use(middlewares.SetDefautlHeaderMiddleware)

	v1 := api.Group("/v1", middlewares.SetHeaderV1Middleware)

	v1.Post("/new-user", handlers_v1.CreateUser)
	v1.Patch("/update-common-data-user", handlers_v1.UpdateCommonDataUser)
	v1.Patch("/update-password-user", handlers_v1.UpdatePasswordUser)
	v1.Delete("/disable-user", handlers_v1.DisableUser)
	v1.Get("/list-users", handlers_v1.ListUsers)
	v1.Get("/user/:username<minLen(3)\\maxLen(150)>", handlers_v1.GetUser)

	log.Fatal(api.Listen(":" + os.Getenv("SERVER_PORT")))
}
