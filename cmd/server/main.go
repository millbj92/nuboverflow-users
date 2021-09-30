package main

import (
	"log"

	"github.com/millbj92/nuboverflow-users/internal/repository"
	"github.com/millbj92/nuboverflow-users/internal/transport/http"
	"github.com/millbj92/nuboverflow-users/internal/user/service"
)

func Run() error {
	userStore, err := repository.New()
	if err != nil {
		return err
	}

	userService := service.New(userStore)
	app, err := http.New(userService)
	if err != nil {
		return err
	}

	if err := app.App.Listen(":3000"); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}