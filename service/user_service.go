package service

import (
	"context"
	"database/sql"
	"google-oauth/helper"
	"google-oauth/model"
	"google-oauth/repository"
	"google-oauth/web"
	"github.com/go-playground/validator/v10"
)

type UserService struct {
	Repository repository.UserRepository
	DB         *sql.DB
	Validate   *validator.Validate
}

func NewUserService(repo repository.UserRepository, db *sql.DB, validator *validator.Validate) *UserService {
	return &UserService{
		Repository: repo,
		DB: db,
		Validate: validator,
	}
}

func (service *UserService) RegisterFromGoogle(ctx context.Context, request model.AuthUser) web.UserResponse {
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

	userRegister := service.Repository.RegisterFromGoogle(ctx, tx, user)	
	return helper.ToUserResponse(userRegister)
}

func (service *UserService) RegisterDefault(ctx context.Context, request web.UserRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		return web.UserResponse{}
	}

	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	user := model.AuthUser{
		Name: request.Name,
		Email: request.Email,
		Password: request.Password,
	}

	userRegister := service.Repository.RegisterDefault(ctx, tx, user)	
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
