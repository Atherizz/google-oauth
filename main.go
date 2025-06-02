package main

import (
	"google-oauth/app"
	"net/http"
)

func main() {
	router := app.NewRouter()

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
