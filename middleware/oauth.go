package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthConfig = oauth2.Config{
	ClientID:     "1035471242348-e1n7ujn46982ibko0s4v3mhf54lbt12n.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8000/callback",
	Scopes:       []string{"openid", "profile", "email"},
	Endpoint:     google.Endpoint,
}

type Oauth2Middleware struct {
	Handler http.Handler
}

func NewOauth2Middleware(handler http.Handler) *Oauth2Middleware {
	return &Oauth2Middleware{
		Handler: handler,
	}
}

func loadTokenFromRequest(request *http.Request) (*oauth2.Token, error) {
	// get data from cookies
	cookie, err := request.Cookie("oauth_token")
	if err != nil {
		return nil, err
	}

	// Mengambil string yang di-encode base64, lalu decode ke bytes.
	tokenBytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	// Mengambil JSON (byte format) lalu mengubahnya ke struct.
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (middleware *Oauth2Middleware) Wrap(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
		token, err := loadTokenFromRequest(request)
		if err != nil || !token.Valid() {
			http.Redirect(writer, request, OauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), http.StatusNotFound)
			return
		}

		// idToken, ok := token.Extra("id_token").(string)

		// if !ok {
		// 	http.Error(writer, "no id_token in field token", http.StatusInternalServerError)
		// }

		// tokenPayload, err := helper.DecodeIdToken(idToken)
		// if err != nil {
		// 	http.Redirect(writer, request, OauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), http.StatusNotFound)
		// 	return
		// }

		// ctx := context.WithValue(request.Context(), "email", tokenPayload.Email)
		// ctx = context.WithValue(ctx, "name", tokenPayload.Name)
		// ctx = context.WithValue(ctx, "picture", tokenPayload.Picture)
	
		tokenSource := OauthConfig.TokenSource(request.Context(), token)

		token, err = tokenSource.Token()
		if err != nil || !token.Valid() {
			http.Redirect(writer, request, OauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), http.StatusNotFound)
			return
		}

		next(writer, request, param)

	}
}
