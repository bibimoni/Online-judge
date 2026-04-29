package routes

import (
	"fmt"
	"problem/storage"
	"problem/utils"
	"problem/utils/polygon"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ProblemRoute(router fiber.Router) {
	router.Post("/add", func(c *fiber.Ctx) error {
		var problemId int
		var err error

		problemId, err = strconv.Atoi(c.Query("problemId", ""))
		if err != nil {
			utils.WriteFailedOutput(c, err)
			return err
		}

		if err := storage.AddProblem(uint64(problemId)); err != nil {
			// return c.Status(500).SendString(fmt.Sprintf("error adding problem: %s", err.Error()))
			utils.WriteFailedOutput(c, err)
			return err
		}

		OK := "OK"
		utils.WriteCreatedOutput(c, &OK, nil)

		return nil
	})

	router.Get("/latest-version", func(c *fiber.Ctx) error {
		var problemId int
		problemId, err := strconv.Atoi(c.Query("problemId", ""))
		if err != nil {
			// return c.Status(500).SendString("something wrong with your problemId parameter")
			utils.WriteFailedOutput(c, err)
		}

		versionNumber, err := polygon.GetLatestVersionNumber(uint64(problemId))
		if err != nil {
			// return c.Status(500).SendString(fmt.Sprintf("error getting latest version number: %s", err.Error()))
			utils.WriteFailedOutput(c, err)
			return err
		}

		versionStr := fmt.Sprintf("v%d", versionNumber)
		utils.WriteSuccessOutput(c, &versionStr, nil)

		return nil
	})

	router.Get("all", func(c *fiber.Ctx) error {
		list, err := storage.GetAllProblems()
		if err != nil {
			// return c.Status(500).SendString(err.Error())
			utils.WriteFailedOutput(c, err)
			return err
		}

		utils.WriteSuccessOutput(c, &list, err)
		return nil
	})

	router.Static("/get/", "/storage")
}
