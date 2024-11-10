package models

import "time"

var schemaArtis = `
CREATE TABLE stream.artis (
	artis_id uuid DEFAULT gen_random_uuid() NOT NULL,
	slug varchar(100) NOT NULL,
	"name" varchar(100) NOT NULL,
	nationality varchar(50) NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp NULL,
	CONSTRAINT artis_pk PRIMARY KEY (artis_id),
	CONSTRAINT artis_unique UNIQUE (artis_slug)
);
`

type Artis struct {
	Artis_Id    string     `db:"artis_id" json:"artis_id,omitempty" form:"artis_id"`
	Slug        string     `db:"slug" json:"slug,omitempty" form:"slug"`
	Name        string     `db:"name" json:"name" form:"name"`
	Nationality string     `db:"nationality" json:"nationality" form:"nationality"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdateAt    *time.Time `db:"updated_at" json:"updated_at"`
}

type Arties []Artis
