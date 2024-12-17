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