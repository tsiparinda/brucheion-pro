 -- create extension IF NOT EXISTS  hstore;
drop table IF EXISTS users;
DROP TABLE IF EXISTS citedata;

CREATE TABLE IF NOT EXISTS public.users (
	id serial4 NOT NULL,
	username text NOT NULL,
	password_hash text NOT NULL,
	is_verified bool NOT NULL DEFAULT false,
	verification_code text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_is_lowercase CHECK ((username = lower(username))),
	CONSTRAINT users_username_key UNIQUE (username)
);

-- public.citedata definition

-- Drop table

-- DROP TABLE public.citedata;

CREATE TABLE public.citedata (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	bucket text NULL,
	boltdb public."hstore" NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT citedata_pkey PRIMARY KEY (id)
);


-- public.citedata foreign keys

ALTER TABLE public.citedata ADD CONSTRAINT citedata_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

CREATE INDEX idx_b_boltdb ON citedata USING BTREE (boltdb);
CREATE INDEX idx_h_bucket ON citedata USING hash (bucket);
CREATE INDEX idx_h_user_id ON citedata USING hash (user_id);