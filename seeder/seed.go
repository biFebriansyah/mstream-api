package main

import (
	"biFebriansyah/gostream/utils"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database := utils.NewDatabase()
	defer database.Close()

	seedGenres(database.DB)
	seedArtis(database.DB)
	seedMusic(database.DB)

	log.Println("seed table succsess")
}
