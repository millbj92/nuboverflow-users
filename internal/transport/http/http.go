package http

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/millbj92/nuboverflow-users/internal/user"
	"github.com/millbj92/nuboverflow-users/internal/user/service"
)

type HttpError struct {
	Message string
	Errors  []error
}

type HealthCheckResponse struct {
	HTTPService string
	Database string
}

type Handler struct {
	App         *fiber.App
	UserService service.UserService
}

func New(UserService service.UserService) (Handler, error) {
	//Init
	app := fiber.New()
	h := Handler{
		App: app,
	}

	//Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept",
	}))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))


	// Version 1
	prefix := app.Group("/api/v1")
	v1 := prefix.Get("/users", h.GetAllUsers)
	v1.Get("users/:id", h.GetUserByID)
	v1.Get("/users", h.GetUserByEmail)
	v1.Post("/users", h.CreateUser)
	v1.Put("/users", h.UpdateUser)
	v1.Delete("/users/:id", h.DeleteUser)
	v1.Get("/ping", h.Healthcheck)
	v1.Get("/dashboard", monitor.New())

	//Ship it
	return h, nil
}

func (h Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		log.Printf("UserService failed to GET /users\nError: %s", err)
		return err
	}
	err = c.JSON(users)
	if err != nil {
		log.Printf("Failed to response to GET /users: %s", err)
		return err
	}
	return nil
}

func (h Handler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.UserService.GetUserByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.Status(fiber.StatusNotFound).JSON(HttpError{
				Message: "Resource was not found.",
				Errors:  []error{err},
			})
		}
		log.Printf("UserService failed to GetUserByID: %s", err)
		return err
	}
	err = c.JSON(result)
	if err != nil {
		log.Printf("Failed to response to GET users/%s\nError: %s", id, err)
		return err
	}
	return nil
}

func (h Handler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
		user, err := h.UserService.GetUserByEmail(email)
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

func (h Handler) CreateUser(c *fiber.Ctx) error {
	usr := new(user.User)
		if err := c.BodyParser(usr); err != nil {
			log.Printf("Error parsing user: %s", err)
			return err
		}
		user, err := h.UserService.CreateUser(*usr)
		if err != nil {
			log.Printf("Error calling CreateUser: %s", err)
			return err
		}
		if err = c.JSON(user); err != nil {
			log.Printf("Error responding to POST /users: %s", err)
			return err
		}
		return nil
}

func (h Handler) UpdateUser(c *fiber.Ctx) error {
	usr := new(user.User)
		if err := c.BodyParser(usr); err != nil {
			log.Printf("Error parsing user: %s", err)
			return err
		}
		user, err := h.UserService.UpdateUser(*usr)
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

func (h Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
		if err := h.UserService.DeleteUser(id); err != nil {
			log.Printf("Error deleting user: %s", err)
			return err
		}
		if err := c.SendStatus(fiber.StatusOK); err != nil {
			log.Printf("Error responding to DELETE /user/:id\nid: %s\nError: %s", id, err)
			return err
		}
		return nil
}

func (h Handler) Healthcheck(c *fiber.Ctx) error {
	err := c.SendString("pong"); if err != nil {
		log.Printf("Error responding to health check: %s", err)
		return err
	}
	return nil
}
