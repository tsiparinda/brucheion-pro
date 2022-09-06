
CREATE TABLE IF NOT EXISTS init (
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
