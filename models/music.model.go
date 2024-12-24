package models

import (
	"time"

	"github.com/lib/pq"
)

var schemaMusic = `
CREATE TABLE stream.music (
	music_id uuid DEFAULT gen_random_uuid() NOT NULL,
	slug varchar(100) NOT NULL,
	title varchar(100) NOT NULL,
	release_date date NOT NULL,
	cover varchar NOT NULL,
	source_url varchar NOT NULL,
	created_at timestamp DEFAULT NOW() NULL,
	updated_at timestamp NULL,
	CONSTRAINT musics_pk PRIMARY KEY (music_id),
	CONSTRAINT musics_unique UNIQUE (slug)
);
`
var indexMusic = `
CREATE UNIQUE INDEX idx_music ON stream.music(slug, music_id);
`

type Music struct {
	Music_id     string         `db:"music_id" json:"music_id,omitempty" form:"music_id"`
	Artis_id     *string        `db:"artis_id" json:"artis_id" form:"artis_id"`
	Title        string         `db:"title" json:"title" form:"title"`
	Slug         string         `db:"slug" json:"slug" form:"slug"`
	Release_date *string        `db:"release_date" json:"release_date" form:"release_date"`
	Cover        string         `db:"cover" json:"cover" form:"cover"`
	Source_url   string         `db:"source_url" json:"source_url" form:"source_url"`
	Genre        pq.StringArray `db:"genres" json:"genres,omitempty" form:"genres"`
	CreatedAt    *time.Time     `db:"created_at" json:"created_at"`
	UpdateAt     *time.Time     `db:"updated_at" json:"updated_at"`
}

type Musics []Music

type MusicData struct {
	Music_id     string       `db:"music_id" json:"music_id,omitempty" form:"music_id"`
	Title        string       `db:"title" json:"title" form:"title"`
	Slug         string       `db:"slug" json:"slug" form:"slug"`
	Release_date string       `db:"release_date" json:"release_date" form:"release_date"`
	Cover        string       `db:"cover" json:"cover" form:"cover"`
	Source_url   string       `db:"source_url" json:"source_url" form:"source_url"`
	MusicArtis   MusicArtis   `json:"MusicArtis" xml:"MusicArtis" form:"MusicArtis"`
	MusicGenre   []GenreMusic `json:"MusicGenre" xml:"MusicGenre" form:"MusicGenre"`
	CreatedAt    *time.Time   `db:"created_at" json:"created_at"`
	UpdateAt     *time.Time   `db:"updated_at" json:"updated_at"`
}

type MusicDatas []MusicData
