package app

import (
	"database/sql"
	"google-oauth/helper"
	"time"
)

func NewDB() *sql.DB {
	dbName := helper.LoadEnv("DB_NAME")
	port := helper.LoadEnv("PORT")
	dbUser := helper.LoadEnv("DB_USER")

	db, err := sql.Open("mysql", dbUser+"@tcp(localhost:"+port+")/"+dbName+"?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
