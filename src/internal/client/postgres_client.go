package client

import (
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

type PostgresClientProvider struct{}

func (p PostgresClientProvider) Connect(env string) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=go_todo_app_"+env+" port=5432 sslmode=disable TimeZone=Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	boil.SetDB(db)
	boil.DebugMode = true
	DB = db
}

func (p PostgresClientProvider) Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}

func (p PostgresClientProvider) Set(db *sql.DB) {
	DB = db
}
