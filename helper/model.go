package helper

import (
	"google-oauth/model"
	"google-oauth/web"
)

func ToUserResponse(user model.AuthUser) web.UserResponse {
	return web.UserResponse{
		Id: user.Id,
		GoogleId: user.GoogleId,
		Name: user.Name,
		Email: user.Email,
		Picture: user.Picture,
		Provider: user.Provider,
		Role: user.Role,
	}
}