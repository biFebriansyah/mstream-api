package main

import (
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"log"

	"github.com/jmoiron/sqlx"
)

func seedMusic(db *sqlx.DB) {
	repos := repositories.NewMusic(db)
	var datas models.MusicDatas = models.MusicDatas{
		{
			Music_id:     "4143e56b-2888-4b87-9436-bf611c5e65dc",
			Title:        "ku katakan dengan indah",
			Slug:         "ku-katakan-dengan-indah",
			Release_date: "2024-02-02",
			Cover:        "https://streamapp.nos.jkt-1.neo.id/9329b9a1.jpg",
			Source_url:   "https://streamapp.nos.jkt-1.neo.id/e3c366d8/master.m3u8",
			MusicArtis:   models.MusicArtis{Artis_id: "f68e6f5f-e9ff-46f9-81a3-4ed4cd2290dd"},
			MusicGenre:   []models.GenreMusic{{Genre_id: "79b529b7-76af-4e80-8eeb-8f017a87b2c1"}},
		},
		{
			Music_id:     "a4c1c0b9-04d2-427d-991d-f358efdd1810",
			Title:        "alamat palsu",
			Slug:         "alamt-palsu",
			Release_date: "2022-02-02",
			Cover:        "https://streamapp.nos.jkt-1.neo.id/9329b9a1.jpg",
			Source_url:   "https://streamapp.nos.jkt-1.neo.id/e3c366d8/master.m3u8",
			MusicArtis:   models.MusicArtis{Artis_id: "80b0bd36-8a66-429f-ab9f-d6b075fe27bc"},
			MusicGenre: []models.GenreMusic{
				{Genre_id: "79b529b7-76af-4e80-8eeb-8f017a87b2c1"},
				{Genre_id: "40019107-d352-4f2e-8619-fbdfe2711a23"},
			},
		},
	}

	if err := repos.InsertBatchData(&datas, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("seed music created")
}
