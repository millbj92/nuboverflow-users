package http

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"
	_ "github.com/millbj92/nuboverflow-users/internal/transport/http/docs"
	"golang.org/x/crypto/bcrypt"

	swagger "github.com/arsmn/fiber-swagger/v2"
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

type CreateUserRequest struct {
	UserName   string `json:"username" validate:"required,min=4,max=100"`
	Password   string `json:"password" validate:"passwd"`
	Email      string `json:"email" validate:"required,email"`
}

// @title Nuboverflow - Users Microservice
// @version 1.0
// @description Used for creation of users within the Nuboverflow domain.

// @contact.name Millbj92]
// @contact.url nuboverflow.com
// @contact.email admin@nuboverflow.com

// @license.name MIT
// @license.url https://github.com/millbj92/nuboverflow-users/blob/main/LICENSE
func CreateRoutes(service usr.Service, v *validator.Validate) *fiber.App {

	registerValidators(v)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/docs/*", swagger.Handler)
	v1 := app.Group("/api/v1")
	v1.Get("/users", GetAllUsers(service))
	v1.Post("/users", CreateUser(service, v))
	v1.Get("/users", GetUserByEmail(service))
	v1.Put("/users", UpdateUser(service))
	v1.Get("users/:id", GetUserByID(service))
	v1.Delete("/users/:id", DeleteUser(service))
	v1.Get("/ping", Healthcheck())
	v1.Get("/dashboard", monitor.New())

	return app
}

// GetAllUsers godoc
// @Summary List all users
// @Description Get all user accounts
// @Tags users
// @Produce  json
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users [get]
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


// GetUserByID godoc
// @Summary Get a single user by their ID
// @Description get user by ID
// @Tags users
// @Accept  int
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/{id} [get]
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


// GetUserByEmail godoc
// @Summary Get a user by their email address.
// @Description get user by ID
// @Tags users
// @Accept  string
// @Produce  json
// @Param q query string false "user search by q" Format(email)
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/{email} [get]
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

// CreateUser godoc
// @Summary Create a user
// @Description add by json user
// @Tags users
// @Accept  json
// @Produce  json
// @Param account body model.CreateUser true "Create user"
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts [post]
func CreateUser(service usr.Service, v *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := CreateUserRequest{}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(503).SendString(err.Error())
		}

		if err := v.Struct(requestBody); err != nil {
			log.Println(err)
			return err
		}

		//Validate if user already exists
		existingUser, err := service.GetUserByEmail(requestBody.Email)
		if err != nil {
			log.Println(err)
		}
		if existingUser.ID > 0 {
			log.Println("Returning error: User exists.")
		    return errors.New("user exists")
		}

		//Hash password.
		passBytes := []byte(requestBody.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		requestBody.Password = string(hashedPassword)

		domainUser := user.User{
			Email: requestBody.Email,
			UserName: requestBody.UserName,
			Password: requestBody.Password,
		}
		//Send to service to be stored in the Store.
		user, err := service.CreateUser(&domainUser)
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

// UpdateUser godoc
// @Summary Update a user
// @Description update by json user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body model.UpdateUser true "Update user"
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts [put]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete by user ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID" Format(int64)
// @Success 204 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/{id} [delete]
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

// Healthcheck godoc
// @Summary Healthcheck the Users API
// @Description Ping this endpoint to get a current healthcheck.
// @Tags users
// @Produce  string
// @Success 200 {object} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/ping [get]
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

func registerValidators(v *validator.Validate) {
	var mustHave = []func(rune) bool{
		unicode.IsUpper,
		unicode.IsLower,
		unicode.IsPunct,
		unicode.IsDigit,
	}
	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {	
		if len(fl.Field().String()) < 8 {
			log.Println("Length < 8??")
			return false
		}
		found := false
		for _, testRune := range mustHave {
			for _, r := range fl.Field().String() {
				if testRune(r) {
					found = true
				}
			}
		}
		return found
	})
}