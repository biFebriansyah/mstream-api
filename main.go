package main

import (
	"biFebriansyah/gostream/routers"
	"biFebriansyah/gostream/utils"
	"context"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	example()
	database := utils.NewDatabase()
	amqpConnection := utils.NewAmqpConn()
	routers := routers.New(database.DB, amqpConnection)
	server := utils.NewServer(routers)

	go amqpConnection.NewConsumer()

	wait := utils.GracefulShutdown(context.Background(), 2*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return database.Close()
		},
		"server": func(ctx context.Context) error {
			return server.Shutdown()
		},
		"Amqp": func(ctx context.Context) error {
			return amqpConnection.Conn.Close()
		},
	})

	<-wait
}
