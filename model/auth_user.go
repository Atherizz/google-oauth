package model

import "time"

type AuthUser struct {
	Id          int
	GoogleId    string
	Name        string
	Email       string
	Picture     string
	Provider    string
	Role        string
	Password 	string
	LastLoginAt time.Time
}
