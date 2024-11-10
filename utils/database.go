package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SqlxDB struct {
	*sqlx.DB
}

func NewDatabase() *SqlxDB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)

	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Fatal("database fail to connect")
	}

	db.SetConnMaxIdleTime(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return &SqlxDB{db}
}

func (s *SqlxDB) Shutdown() error {
	return s.Close()
}
