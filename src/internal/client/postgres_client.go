package client

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type PostgresClientProvider struct {
	sqlDB *sql.DB
}

func (p PostgresClientProvider) Connect(env string) {
	dsn := "host=0.0.0.0 user=postgres password=postgres dbname=go_todo_app_" + env + " port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(2)

	DB = db
	p.sqlDB = sqlDB
}

func (p PostgresClientProvider) Close() {
	if err := p.sqlDB.Close(); err != nil {
		panic(err)
	}
}
