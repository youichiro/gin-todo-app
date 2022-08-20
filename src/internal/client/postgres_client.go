package client

import (
	"database/sql"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvidePostgreSqlClient() (*sql.DB, *gorm.DB) {
	url := os.Getenv("POSTGRESQL_URL")
	log.Print(url)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(2)

	return sqlDB, db
}
