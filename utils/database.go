package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SqlxDB struct {
	*sqlx.DB
}

var (
	instance *SqlxDB
	once     sync.Once
)

func NewDatabase() *SqlxDB {
	once.Do(func() {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")

		config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)

		db, err := sqlx.Connect("postgres", config)
		if err != nil {
			log.Fatalf("database fail to connect with config %s", config)
		}

		db.SetConnMaxIdleTime(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)

		instance = &SqlxDB{db}
	})
	return instance
}

func (s *SqlxDB) Shutdown() error {
	return s.Close()
}
