package client

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

type PostgresClientProvider struct{}

func (p PostgresClientProvider) Connect(env string) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := "go_todo_app_" + env
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Asia/Tokyo"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	DB = db
}

func (p PostgresClientProvider) Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}
