package repository

import (
	"context"
	"database/sql"
	"errors"
	"google-oauth/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) RegisterFromGoogle(ctx context.Context, tx *sql.Tx, user model.AuthUser) model.AuthUser {
	script := "INSERT INTO users(google_id,name,email,picture,provider,role) VALUES (?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, script, user.GoogleId, user.Name, user.Email, user.Picture, user.Provider, user.Role)
	if err != nil {
		return model.AuthUser{}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.AuthUser{}
	}
	user.Id = int(id)
	return user

}

func (repo *UserRepository) RegisterDefault(ctx context.Context, tx *sql.Tx, user model.AuthUser) model.AuthUser {
	script := "INSERT INTO users(name,email,password) VALUES (?,?,?)"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashedString := string(hashedPassword)

	result, err := tx.ExecContext(ctx, script,user.Name, user.Email, hashedString)
	if err != nil {
		return model.AuthUser{}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.AuthUser{}
	}
	user.Id = int(id)
	return user

}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, tx *sql.Tx, email string) (model.AuthUser, error) {
	script := "SELECT id,google_id,name,email,picture,provider,role FROM users WHERE email = (?)"
	result, err := tx.QueryContext(ctx, script, email)
	if err != nil {
		return model.AuthUser{}, err
	}

	user := model.AuthUser{}

	if result.Next() {
		err := result.Scan(&user.Id, &user.GoogleId, &user.Name, &user.Email, &user.Picture, &user.Provider, &user.Role)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	return user, errors.New("ID not found")

}
