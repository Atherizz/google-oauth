package controller

import (
	"fmt"
	"google-oauth/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func BasicOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Fprint(writer, "selamat datang di endpoint basic auth! anda berhasil terautentikasi \n")
}

func HomeOauth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Fprint(writer, "welcome, you are authenticated \n")
}

func Callback(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	code := request.URL.Query().Get("code")
	token, err := middleware.OauthConfig.Exchange(request.Context(), code)
	if err != nil {
		http.Error(writer, "failed get token", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, "you are authenticated! Token", token.AccessToken)
}