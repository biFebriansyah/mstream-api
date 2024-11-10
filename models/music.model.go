package models

import "time"

var schemaMusic = `
CREATE TABLE stream.musics (
	music_id uuid DEFAULT gen_random_uuid() NOT NULL,
	artis uuid NULL,
	title varchar(100) NOT NULL,
	release_date date NOT NULL,
	cover varchar NOT NULL,
	source_url varchar NOT NULL,
	created_at timestamp DEFAULT NOW() NULL,
	updated_at timestamp NULL,
	CONSTRAINT musics_pk PRIMARY KEY (music_id),
	CONSTRAINT musics_artis_fk FOREIGN KEY (artis) REFERENCES stream.artis(artis_id) ON DELETE SET NULL
);
`

type Music struct {
	Music_id     string     `db:"music_id" json:"music_id,omitempty" form:"music_id"`
	Artis        string     `db:"artis" json:"artis" form:"artis"`
	Title        string     `db:"title" json:"title" form:"title"`
	Release_date string     `db:"release_date" json:"release_date" form:"release_date"`
	Cover        string     `db:"cover" json:"cover" form:"cover"`
	Source_url   string     `db:"source_url" json:"source_url" form:"source_url"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdateAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type Musics []Music
