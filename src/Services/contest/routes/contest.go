package routes

import (
	repository "contest/domain/repository/contest/impl"
	"contest/internal/infrastructure/database"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ContestRoutes(app *fiber.App) {
	route := app.Group("/contest")

	route.Post("/create", func(c *fiber.Ctx) error {
		cr := repository.NewContestRepository(database.Db)

		var cancel context.CancelFunc
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var authorId uint64
		var err error
		if authorId, err = strconv.ParseUint(c.Query("author-id", "0"), 10, 64); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		if authorId == 0 {
			return c.Status(400).SendString("author-id is required")
		}

		contestId, err := cr.Create(ctx, authorId)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(map[string]interface{}{
			"contest-id": contestId,
		})
	})

	route.Post("/edit", func(c *fiber.Ctx) error {
		var editType string

		if editType = c.Query("edit-type", ""); editType == "" {
			return c.Status(400).SendString("contest-id required")
		}

		if editType == "add-people" {
			return handleAddPeople(c)
		}

		if editType == "remove-people" {
			return handleRemovePeople(c)
		}

		if editType == "edit-problemset" {
			return handleEditProblemset(c)
		}

		return c.Status(400).SendString("invalid edit-type")
	})
}

func handleAddPeople(c *fiber.Ctx) error {
	type AddPeopleRequest struct {
		ContestId  string `json:"contest-id" validate:"required"`
		PeopleType string `json:"people-type" validate:"required"`
		UserId     uint64 `json:"user-id" validate:"min=1"`
	}

	req := new(AddPeopleRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	cr := repository.NewContestRepository(database.Db)
	if err := cr.AddPeople(req.ContestId, req.PeopleType, req.UserId); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func handleRemovePeople(c *fiber.Ctx) error {
	type RemovePeopleRequest struct {
		ContestId  string `json:"contest-id" validate:"required"`
		PeopleType string `json:"people-type" validate:"required"`
		UserId     uint64 `json:"user-id" validate:"min=1"`
	}

	req := new(RemovePeopleRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	cr := repository.NewContestRepository(database.Db)
	if err := cr.RemovePeople(req.ContestId, req.PeopleType, req.UserId); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func handleEditProblemset(c *fiber.Ctx) error {
	type EditProblemsetRequest struct {
		ContestId     string   `json:"contest-id" validate:"required"`
		NewProblemset []uint64 `json:"new-problemset" validate:"required"`
	}

	req := new(EditProblemsetRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	cr := repository.NewContestRepository(database.Db)
	if err := cr.EditProblemset(req.ContestId, req.NewProblemset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}
