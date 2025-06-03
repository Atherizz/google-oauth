package middleware

import (
	"context"
	"google-oauth/helper"
	"google-oauth/model"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: handler,
	}
}

func (middleware *AuthMiddleware) Wrap(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
		session, _ := helper.Store.Get(request, "user_info")

		user, ok := session.Values["user"].(model.AuthUser)
		if !ok ||  user.Name == "" || user.Email == "" {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "user", user)

		next(writer, request.WithContext(ctx), param)
	}
}
