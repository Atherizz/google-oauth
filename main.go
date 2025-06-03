package main

import (
	"encoding/gob"
	"google-oauth/app"
	"google-oauth/model"
	"net/http"
)

func main() {
	router := app.NewRouter()
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
