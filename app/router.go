package app

import (
	"google-oauth/handler"
	"google-oauth/middleware"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	oauthMiddleware := middleware.NewOauth2Middleware(router)
	authMiddleware := middleware.NewAuthMiddleware(router)

	router.GET("/auth/google/login", handler.LoginOauth)

	router.GET("/login", handler.LoginView)
	router.GET("/home", oauthMiddleware.Wrap(handler.HomeOauth))
	router.GET("/callback", handler.Callback)
	router.GET("/api/user", oauthMiddleware.Wrap(authMiddleware.Wrap(handler.ProfileApi)))
	router.GET("/profile", oauthMiddleware.Wrap(authMiddleware.Wrap(handler.ProfileOauth)))
	router.GET("/logout", oauthMiddleware.Wrap(handler.Logout))

	return router

}
