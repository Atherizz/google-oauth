package app

import (
	"google-oauth/controller"
	"google-oauth/middleware"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	oauthMiddleware := middleware.NewOauth2Middleware(router)
	router.GET("/home", oauthMiddleware.Wrap(controller.HomeOauth))
	router.GET("/callback", controller.Callback)
	router.GET("/api/user", oauthMiddleware.Wrap(controller.GetUser))
	router.GET("/logout", oauthMiddleware.Wrap(controller.Logout))

	return router

}
