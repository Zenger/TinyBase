package main

import (
	"TinyBase/config"
	"TinyBase/db"
	"TinyBase/handlers"
	"TinyBase/shared"

	"fmt"
	"github.com/gofiber/fiber/v3"
	"log"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		panic(err)
	}

	database, err := db.Connect(settings.Database.Host, settings.Database.Port, settings.Database.Username, settings.Database.Password, settings.Database.Database)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	err = db.Bootstrap(database, settings.App.SuperUser, settings.App.Salt)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Get("/:name/:age/:gender?", func(c fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is a %s at %s years old", c.Params("name"), c.Params("gender"), c.Params("age"))
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	tinyBaseContext := shared.TinyBaseContext{
		Database: database,
		Settings: settings,
	}

	app.Post("/sign-in", func(c fiber.Ctx) error {
		return handlers.AuthHandler(c, tinyBaseContext)
	})

	// host := fmt.Sprintf(":%s", settings.App.Port)
	log.Fatal(app.Listen(":6722"))
	defer database.Close()

}
