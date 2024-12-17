CREATE TABLE stream.music_genre (
	music_id uuid NULL,
	genre_id uuid NULL,
	CONSTRAINT music_genre_fk1 FOREIGN KEY (genre_id) REFERENCES stream.genre(genre_id) ON DELETE CASCADE,
	CONSTRAINT music_genre_fk2 FOREIGN KEY (music_id) REFERENCES stream.music(music_id) ON DELETE CASCADE,
	CONSTRAINT music_genre_unique UNIQUE (music_id, genre_id)
);