package handler

import (
	"google-oauth/helper"
	"google-oauth/model"
	"google-oauth/web"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ProfileApi(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	user := request.Context().Value("user")

	authUser := user.(model.AuthUser)

	userResponse := web.UserResponse{
		Email:   authUser.Email,
		Name:    authUser.Name,
		Picture: authUser.Picture,
	}

	helper.WriteEncodeResponse(writer, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	})
}