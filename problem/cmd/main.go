package main

import (
	"fmt"
	"log"
	"os"
	"problem/routes"
	"problem/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New(fiber.Config{
		ServerHeader: "HCMUT-OJ",
	})

	if err := storage.GetMongoDbClient(); err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request!")
	})

	routes.ProblemRoute(app.Group("/problem"))

	for _, route := range app.Stack() {
		for _, r := range route {
			fmt.Printf("%-6s %s\n", r.Method, r.Path)
		}
	}

	port := os.Getenv("PROBLEM_PORT")
	if port == "" {
		log.Fatalln("Can't get port address")
		return
	}

	host := os.Getenv("PROBLEM_HOST")
	if host == "" {
		log.Fatalln("Can't get host address")
		return
	}

	databaseName := os.Getenv("PROBLEM_DATABASE_NAME")
	if databaseName == "" {
		log.Fatalln("Can't get database name")
		return
	}

	serverAddr := host + ":" + port

	log.Printf("Problem server is listening on: %s", serverAddr)
	app.Listen(serverAddr)
}
