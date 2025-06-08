package main

import (
	"encoding/gob"
	"google-oauth/app"
	"google-oauth/handler"
	"google-oauth/model"
	"google-oauth/repository"
	"google-oauth/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()

	validate := validator.New(validator.WithRequiredStructEnabled())

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(*userRepository,db, validate)
	userController := handler.NewOauthController(userService)

	router := app.NewRouter(userController)
	gob.Register(model.AuthUser{})

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
