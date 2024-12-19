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
			Music_id:     "61cbea2c-5af2-4733-b313-b057e8da1b35",
			Title:        "ku katakan dengan indah",
			Slug:         "ku-katakan-dengan-indah",
			Release_date: "2024-02-02",
			Cover:        "https://dummyimage.com/600x400/cfcfcf/fff",
			Source_url:   "https://dummyimage.com/600x400/cfcfcf/fff",
			MusicArtis:   models.MusicArtis{Artis_id: "f68e6f5f-e9ff-46f9-81a3-4ed4cd2290dd"},
			MusicGenre:   []models.GenreMusic{{Genre_id: "79b529b7-76af-4e80-8eeb-8f017a87b2c1"}},
		},
		{
			Music_id:     "a4c1c0b9-04d2-427d-991d-f358efdd1810",
			Title:        "ku katakan dengan indah2",
			Slug:         "ku-katakan-dengan-indah2",
			Release_date: "2024-02-02",
			Cover:        "https://dummyimage.com/600x400/cfcfcf/fff",
			Source_url:   "https://dummyimage.com/600x400/cfcfcf/fff",
			MusicArtis:   models.MusicArtis{Artis_id: "f68e6f5f-e9ff-46f9-81a3-4ed4cd2290dd"},
			MusicGenre:   []models.GenreMusic{{Genre_id: "79b529b7-76af-4e80-8eeb-8f017a87b2c1"}},
		},
	}

	for _, v := range datas {
		if _, err := repos.InsertData(&v); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("seed music created")
}
