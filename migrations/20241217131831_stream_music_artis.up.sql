CREATE TABLE stream.music_artis (
	music_id uuid NULL,
	artis_id uuid NULL,
	CONSTRAINT music_artis_fk1 FOREIGN KEY (music_id) REFERENCES stream.music(music_id) ON DELETE CASCADE,
	CONSTRAINT music_artis_fk2 FOREIGN KEY (artis_id) REFERENCES stream.artis(artis_id) ON DELETE CASCADE,
	CONSTRAINT music_artis_unique1 UNIQUE (music_id, artis_id)
);