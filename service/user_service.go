package service

import (
	"context"
	"database/sql"
	"google-oauth/helper"
	"google-oauth/model"
	"google-oauth/repository"
	"google-oauth/web"
)

type UserService struct {
	Repository repository.UserRepository
	DB         *sql.DB
}

func NewUserService(repo repository.UserRepository, db *sql.DB) *UserService {
	return &UserService{
		Repository: repo,
		DB: db,
	}
}

func (service *UserService) Register(ctx context.Context, request model.AuthUser) web.UserResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	user := model.AuthUser{
		// id,google_id,name,email,picture,provider,role
		Id: request.Id,
		GoogleId: request.GoogleId,
		Name: request.Name,
		Email: request.Email,
		Picture: request.Picture,
		Provider: request.Provider,
		Role: request.Role,
	}

	userRegister := service.Repository.Register(ctx, tx, user)	
	return helper.ToUserResponse(userRegister)
}

func (service *UserService) GetUserByEmail(ctx context.Context, email string) web.UserResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	user,err := service.Repository.GetUserByEmail(ctx, tx, email)	
	if err != nil {
		return web.UserResponse{}
	}

	return helper.ToUserResponse(user)
}
