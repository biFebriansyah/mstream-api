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