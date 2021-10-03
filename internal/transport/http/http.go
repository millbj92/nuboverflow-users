package http

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/millbj92/nuboverflow-users/internal/user"
	usr "github.com/millbj92/nuboverflow-users/internal/user/service"
)

type HttpError struct {
	Message string
	Errors  []error
}

type HealthCheckResponse struct {
	HTTPService string
	Database    string
}

func CreateRoutes(service usr.Service) *fiber.App {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	v1 := app.Group("/api/v1")
	v1.Get("/users", GetAllUsers(service))
	v1.Post("/users", CreateUser(service))
	v1.Get("/users", GetUserByEmail(service))
	v1.Put("/users", UpdateUser(service))
	v1.Get("users/:id", GetUserByID(service))
	v1.Delete("/users/:id", DeleteUser(service))
	v1.Get("/ping", Healthcheck())
	v1.Get("/dashboard", monitor.New())

	return app
}

func GetAllUsers(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := service.GetAllUsers()
		if err != nil {
			log.Printf("UserService failed to GET /users\nError: %s", err)
			return err
		}
		err = c.JSON(&users)
		if err != nil {
			log.Printf("Failed to response to GET /users: %s", err)
			return err
		}
		return nil
	}
}

func GetUserByID(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := intFromString(utils.ImmutableString(c.Params("id")))
		if err != nil {
			return err
		}
		result, err := service.GetUserByID(id)
		if err != nil {
			if err.Error() == "record not found" {
				if err := c.Status(fiber.StatusNotFound).JSON(HttpError{
					Message: "Resource was not found.",
					Errors:  []error{err},
				}); err != nil {
					return err
				}
			}
			log.Printf("UserService failed to GetUserByID: %s", err)
			return err
		}
		err = c.JSON(result)
		if err != nil {
			log.Printf("Failed to response to GET users/%s\nError: %s", fmt.Sprint(id), err)
			return err
		}
		return nil
	}
}

func GetUserByEmail(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email")
		user, err := service.GetUserByEmail(email)
		if err != nil {
			log.Printf("Error calling GetUserByEmail: %s", err)
			return err
		}
		err = c.JSON(user)
		if err != nil {
			log.Printf("Failed to respond to POST /users/email\nEmail: %s\nError: %s", email, err)
			return err
		}
		return nil
	}
}

func CreateUser(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody user.User
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(503).SendString(err.Error())
		}
		user, err := service.CreateUser(&requestBody)
		if err != nil {
			if err.Error() == "user exists" {
				return c.Status(fiber.StatusBadRequest).SendString("User Already Exists")

			} else {
				log.Printf("Error calling CreateUser: %s", err)
				return err
			}
		}
		if err = c.JSON(user); err != nil {
			log.Printf("Error responding to POST /users: %s", err)
			return err
		}
		return nil
	}
}

func UpdateUser(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usr := new(user.User)
		if err := c.BodyParser(usr); err != nil {
			log.Printf("Error parsing user: %s", err)
			return err
		}
		user, err := service.UpdateUser(*usr)
		if err != nil {
			log.Printf("Error calling UpdateUser %s", err)
			return err
		}
		if err = c.JSON(user); err != nil {
			log.Printf("Error responding to PUT /users: %s", err)
			return err
		}
		return nil
	}
}

func DeleteUser(service usr.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := intFromString(utils.ImmutableString(c.Params("id")))
		if err != nil {
			return err
		}
		if err := service.DeleteUser(id); err != nil {
			log.Printf("Error deleting user: %s", err)
			return err
		}
		if err := c.SendStatus(fiber.StatusOK); err != nil {
			log.Printf("Error responding to DELETE /user/:id\nid: %s\nError: %s", fmt.Sprint(id), err)
			return err
		}
		return nil
	}
}

func Healthcheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.SendString("pong")
		if err != nil {
			log.Printf("Error responding to health check: %s", err)
			return err
		}
		return nil
	}
}

func intFromString(sid string) (int, error) {
	uid, err := strconv.ParseInt(sid, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(uid), nil
}
