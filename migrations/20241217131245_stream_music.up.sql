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