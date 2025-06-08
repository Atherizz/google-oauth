package app

import (
	"google-oauth/handler"
	"google-oauth/middleware"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(oauthController *handler.OauthController) *httprouter.Router {
	router := httprouter.New()

	oauthMiddleware := middleware.NewOauth2Middleware(router)
	authMiddleware := middleware.NewAuthMiddleware(router)

	router.GET("/auth/google/login", oauthController.LoginOauth)
	router.GET("/home", oauthMiddleware.Wrap(oauthController.HomeOauth))
	router.GET("/callback", oauthController.Callback)
	router.GET("/profile", oauthMiddleware.Wrap(authMiddleware.Wrap(oauthController.ProfileOauth)))
	router.GET("/logout", oauthMiddleware.Wrap(oauthController.Logout))
	
	router.GET("/login", handler.LoginView)
	router.GET("/register",handler.RegisterView)

	router.GET("/api/user", oauthMiddleware.Wrap(authMiddleware.Wrap(handler.ProfileApi)))
	router.POST("/api/register", oauthController.RegisterDefault)

	return router

}
