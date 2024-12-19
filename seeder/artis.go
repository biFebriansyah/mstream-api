package main

import (
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"log"

	"github.com/jmoiron/sqlx"
)

func seedArtis(db *sqlx.DB) {
	repos := repositories.NewArtis(db)
	dob := "1990-02-20"
	var datas models.Arties = models.Arties{
		{
			Artis_Id:      "f68e6f5f-e9ff-46f9-81a3-4ed4cd2290dd",
			FirstName:     "aril",
			LastName:      "noah",
			Email:         "arilnoah@email.com",
			BirthDate:     &dob,
			Picture:       "https://pinnacle.works/wp-content/uploads/2022/06/dummy-image.jpg",
			Slug:          "arilnoah",
			Nationality:   "indonesia",
			StreetAddress: "jln.satu",
			City:          "jakarta",
			Province:      "jakarta",
			ZipCode:       1234,
		},
		{
			Artis_Id:      "80b0bd36-8a66-429f-ab9f-d6b075fe27bc",
			FirstName:     "ayu",
			LastName:      "ting-ting",
			Email:         "ayutingting@email.com",
			BirthDate:     &dob,
			Picture:       "https://pinnacle.works/wp-content/uploads/2022/06/dummy-image.jpg",
			Slug:          "ayutingting",
			Nationality:   "indonesia",
			StreetAddress: "jln.dua",
			City:          "depok",
			Province:      "depok",
			ZipCode:       1233,
		},
		{
			Artis_Id:      "ade01190-9353-4808-80f0-18eaf316f7bb",
			FirstName:     "iwan",
			LastName:      "false",
			Email:         "iwanfalse@email.com",
			BirthDate:     &dob,
			Picture:       "https://pinnacle.works/wp-content/uploads/2022/06/dummy-image.jpg",
			Slug:          "iwanfalse",
			Nationality:   "indonesia",
			StreetAddress: "jln.tiga",
			City:          "jakarta",
			Province:      "jakarta",
			ZipCode:       1234,
		},
	}

	for _, v := range datas {
		if _, err := repos.InsertData(&v); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("seed artis created")
}
