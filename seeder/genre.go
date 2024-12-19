package main

import (
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"log"

	"github.com/jmoiron/sqlx"
)

func seedGenres(db *sqlx.DB) {
	repos := repositories.NewGenre(db)

	var datas models.Genres = models.Genres{
		{Genre_id: "79b529b7-76af-4e80-8eeb-8f017a87b2c1", Name: "Pop", Slug: "pop"},
		{Genre_id: "fe643ac7-23ce-46fa-80d8-d5e3de38305b", Name: "Rock", Slug: "rock"},
		{Genre_id: "40019107-d352-4f2e-8619-fbdfe2711a23", Name: "Dangdut", Slug: "dangdut"},
	}

	for _, v := range datas {
		if _, err := repos.InsertData(&v); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("seed genres created")
}
