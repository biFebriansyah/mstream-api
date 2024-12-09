package models

import "time"

var schemaGenre = `
CREATE TABLE stream.genre (
	genre_id uuid DEFAULT gen_random_uuid() NOT NULL,
	name varchar NOT NULL,
	slug varchar(100) NOT NULL,
	created_at timestamp DEFAULT NOW() NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT genre_pk PRIMARY KEY (genre_id),
	CONSTRAINT genre_unique UNIQUE (slug)
);
`

var schemaMusicGen = `
CREATE TABLE stream.music_genre (
	music_id uuid NULL,
	genre_id uuid NULL,
	CONSTRAINT music_genre_fk1 FOREIGN KEY (genre_id) REFERENCES stream.genre(genre_id) ON DELETE CASCADE,
	CONSTRAINT music_genre_fk2 FOREIGN KEY (music_id) REFERENCES stream.music(music_id) ON DELETE CASCADE,
	CONSTRAINT music_genre_unique UNIQUE (music_id, genre_id)
);
`
var indexMusicGenres = `CREATE INDEX idx_music_genre ON stream.music_genre(music_id)`

type Genre struct {
	Genre_id  string     `db:"genre_id" json:"genre_id,omitempty" form:"genre_id"`
	Name      string     `db:"name" json:"name" form:"name"`
	Slug      string     `db:"slug" json:"slug" form:"slug"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdateAt  *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type GenreMusic struct {
	Music_id  *string `db:"music_id" json:"music_id,omitempty" form:"music_id"`
	Genre_id  *string `db:"genre_id" json:"genre_id,omitempty" form:"genre_id"`
	GenreName string  `json:"name"`
}

type Genres []Genre
