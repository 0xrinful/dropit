CREATE TYPE user_role AS ENUM('user', 'admin');

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  role user_role NOT NULL DEFAULT 'user'
);
