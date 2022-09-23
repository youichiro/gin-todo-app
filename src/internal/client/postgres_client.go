package client

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type PostgresClientProvider struct{}

func (p PostgresClientProvider) Connect(env string) {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=go_todo_app_"+env+" port=5432 sslmode=disable TimeZone=Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)

	DB = db
}

func (p PostgresClientProvider) Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}
