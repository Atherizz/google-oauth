package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google-oauth/helper"
	"google-oauth/middleware"
	"google-oauth/model"
	"google-oauth/service"
	"google-oauth/web"
	"html/template"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

type OauthController struct {
	Service service.UserService
}

func NewOauthController(service *service.UserService) *OauthController {
	return &OauthController{
		Service: *service,
	}
}

func (controller *OauthController) BasicOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Fprint(writer, "selamat datang di endpoint basic auth! anda berhasil terautentikasi \n")
}

func (controller *OauthController) HomeOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	session, _ := helper.Store.Get(request, "user_info")

	user, ok := session.Values["user"].(model.AuthUser)
	if !ok || user.Email == "" || user.Name == "" {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}
	// fmt.Fprint(writer, "welcome ", name)
	tmpl := template.Must(template.ParseFiles("./resources/welcome.gohtml"))
	tmpl.ExecuteTemplate(writer, "welcome.gohtml", user.Name)

}

func (controller *OauthController) LoginOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	url := middleware.OauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(writer, request, url, http.StatusSeeOther)

}

func (controller *OauthController) RegisterDefault(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	registeredUser := web.UserRequest{}
	err := decoder.Decode(&registeredUser)
	if err != nil {
		panic(err)
	}
	
	response := controller.Service.RegisterDefault(request.Context(), registeredUser)

	helper.WriteEncodeResponse(writer, web.WebResponse{
		Code: 200,
		Status: "OK",
		Data: response,
	})

}

func (controller *OauthController) ProfileOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	session, _ := helper.Store.Get(request, "user_info")

	user, ok := session.Values["user"].(model.AuthUser)
	if !ok || user.Email == "" || user.Name == "" {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}
	// fmt.Fprint(writer, "welcome ", name)
	http.ServeFile(writer, request, "./resources/profile.html")

}

func (controller *OauthController) Callback(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	userResponse := controller.Service.GetUserByEmail(request.Context(), tokenPayload.Email)

	if userResponse.Email == "" {
		userRequest := model.AuthUser{
			GoogleId: tokenPayload.Sub,
			Name:     tokenPayload.Name,
			Email:    tokenPayload.Email,
			Picture:  tokenPayload.Picture,
		}

		controller.Service.RegisterFromGoogle(request.Context(), userRequest)
	}

	session, _ := helper.Store.Get(request, "user_info")
	session.Values["user"] = model.AuthUser{
		Name:     tokenPayload.Name,
		Email:    tokenPayload.Email,
		Picture:  tokenPayload.Picture,
		GoogleId: tokenPayload.Sub,
	}

	err = session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/home", http.StatusFound)
}

func (controller *OauthController) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

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
	http.Redirect(writer, request, "/login", http.StatusFound)
}
