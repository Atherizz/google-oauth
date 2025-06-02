package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google-oauth/helper"
	"google-oauth/middleware"
	"google-oauth/web"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

var store = sessions.NewCookieStore([]byte(helper.LoadEnv("SESSION_SECRET")))

func BasicOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Fprint(writer, "selamat datang di endpoint basic auth! anda berhasil terautentikasi \n")
}

func HomeOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	session, _ := store.Get(request, "user_info")

	name, ok := session.Values["name"].(string)
	if !ok || name == "" {
	http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(writer, "welcome ", name)
	// http.ServeFile(writer, request, "./resources/welcome.html")
}

func Callback(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	code := request.URL.Query().Get("code")
	token, err := middleware.OauthConfig.Exchange(request.Context(), code)
	if err != nil {
		http.Error(writer, "failed get token", http.StatusInternalServerError)
		return
	}

	idToken, ok := token.Extra("id_token").(string)

	if !ok {
		http.Error(writer, "no id_token in field token", http.StatusInternalServerError)
	}

	tokenPayload, err := helper.DecodeIdToken(idToken)
	if err != nil {
		http.Error(writer, "failed decode token", http.StatusInternalServerError)
	}

	// fmt.Fprint(writer, "you are authenticated! Token : \n", token.AccessToken)
	// fmt.Fprint(writer, "\nID Token : \n", idToken)

	tokenJson, err := json.Marshal(token)
	if err != nil {
		http.Error(writer, "failed to marshal token", http.StatusInternalServerError)
		return
	}
	encoded := base64.StdEncoding.EncodeToString(tokenJson)

	http.SetCookie(writer, &http.Cookie{
		Name:     "oauth_token",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	session, _ := store.Get(request, "user_info")
	session.Values["name"] = tokenPayload.Name
	session.Values["email"] = tokenPayload.Email
	session.Values["picture"] = tokenPayload.Picture

	err = sessions.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/home", http.StatusFound)
}

func Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	cookie := http.Cookie{
		Name:     "oauth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(writer, &cookie)
	http.Redirect(writer, request, "/home", http.StatusFound)
}

func GetUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email := request.Context().Value("email").(string)
	name := request.Context().Value("name").(string)
	picture := request.Context().Value("picture").(string)
	if email == "" || name == "" {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userResponse := web.UserResponse{
		Email:   email,
		Name:    name,
		Picture: picture,
	}

	helper.WriteEncodeResponse(writer, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	})
}
