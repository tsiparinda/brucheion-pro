create extension IF NOT EXISTS  hstore;

CREATE TABLE citedata (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	bucket text NULL,
	dict "hstore" NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT citedata_pkey PRIMARY KEY (id)
);

-- citedata foreign keys

ALTER TABLE citedata ADD CONSTRAINT citedata_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_h_dict ON citedata USING hash (dict);
CREATE INDEX idx_h_bucket ON citedata USING hash (bucket);
CREATE INDEX idx_h_user_id ON citedata USING hash (user_id);