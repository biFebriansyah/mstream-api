package models

import (
	"time"
)

var schemaArtis = `
CREATE TABLE stream.artis (
	artis_id uuid DEFAULT gen_random_uuid() NOT NULL,
	slug varchar(100) NOT NULL,
	first_name varchar(50) NOT NULL,
	last_name varchar(50) NOT NULL,
	email varchar(50) NOT NULL,
	birth_date date NULL,
	picture varchar NOT NULL,
	nationality varchar(50) NULL,
	street_address varchar(50) NULL,
	city varchar(50) NULL,
	province varchar(50) NULL,
	zip_code integer NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT artis_pk PRIMARY KEY (artis_id),
	CONSTRAINT artis_unique UNIQUE (slug)
);
`

var schemaArtisMusic = `
CREATE TABLE stream.music_artis (
	music_id uuid NULL,
	artis_id uuid NULL,
	CONSTRAINT music_artis_fk1 FOREIGN KEY (music_id) REFERENCES stream.music(music_id) ON DELETE CASCADE,
	CONSTRAINT music_artis_fk2 FOREIGN KEY (artis_id) REFERENCES stream.artis(artis_id) ON DELETE CASCADE,
	CONSTRAINT music_artis_unique1 UNIQUE (music_id, artis_id)
);
`

var indexMusicArtis = `
CREATE INDEX idx_music_artis ON stream.music_artis(music_id);
CREATE UNIQUE INDEX idx_artis ON stream.artis(artis_id, slug);
`

type Artis struct {
	Artis_Id      string     `db:"artis_id" json:"artis_id,omitempty" form:"artis_id"`
	FirstName     string     `db:"first_name" json:"first_name" form:"first_name"`
	LastName      string     `db:"last_name" json:"last_name" form:"last_name"`
	Email         string     `db:"email" json:"email" form:"email"`
	BirthDate     *string    `db:"birth_date" json:"birth_date" form:"birth_date"`
	Picture       string     `db:"picture" json:"picture" form:"picture"`
	Slug          string     `db:"slug" json:"slug,omitempty" form:"slug"`
	Nationality   string     `db:"nationality" json:"nationality" form:"nationality"`
	StreetAddress string     `db:"street_address" json:"street_address" form:"street_address"`
	City          string     `db:"city" json:"city" form:"city"`
	Province      string     `db:"province" json:"province" form:"province"`
	ZipCode       int        `db:"zip_code" json:"zip_code" form:"zip_code"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdateAt      *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

type MusicArtis struct {
	Music_id  *string `db:"music_id" json:"music_id,omitempty" form:"music_id"`
	Artis_id  string  `db:"artis_id" json:"artis_id,omitempty" form:"artis_id"`
	ArtisName string  `json:"artis_name"`
}

type Arties []Artis
