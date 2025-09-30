package routes

import (
	"fmt"
	"problem/storage"
	"problem/utils/polygon"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func ProblemRoute(router fiber.Router) {
	router.Post("/add", func(c *fiber.Ctx) error {
		var problemId int
		var err error

		problemId, err = strconv.Atoi(c.Query("problemId", ""))
		if err != nil {
			return err
		}

		if err := storage.AddProblem(uint64(problemId)); err != nil {
			return c.Status(500).SendString(fmt.Sprintf("error adding problem: %s", err.Error()))
		}

		return c.SendStatus(200)
	})

	router.Get("/latest", func(c *fiber.Ctx) error {
		var problemId int
		problemId, err := strconv.Atoi(c.Query("problemId", ""))
		if err != nil {
			return c.Status(500).SendString("something wrong with your problemId parameter")
		}

		packageId, err := polygon.GetLastestPackage(uint64(problemId))
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("error getting latest package: %s", err.Error()))
		}

		return c.SendString(strconv.FormatInt(int64(packageId), 10))
	})

	router.Get("all", func(c *fiber.Ctx) error {
		list, err := storage.GetAllProblems()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(bson.M{
			"list": list,
		})
	})

	router.Static("/get/", "/storage")
}
