package models

import "time"

var schemaGenre = `
CREATE TABLE stream.genre (
	genre_id uuid DEFAULT gen_random_uuid() NOT NULL,
	genre_name varchar NOT NULL,
	created_at timestamp DEFAULT NOW() NULL,
	updated_at timestamp NULL,
	CONSTRAINT genre_pk PRIMARY KEY (genre_id)
);
`

type Genre struct {
	Genre_id   string     `db:"genre_id" json:"genre_id,omitempty" form:"genre_id"`
	Genre_name string     `db:"genre_name" json:"genre_name" form:"genre_name"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdateAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type Genres []Genre
