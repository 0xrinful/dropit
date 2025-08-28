CREATE TABLE IF NOT EXISTS files (
  id bigserial PRIMARY KEY,
  token char(8) UNIQUE NOT NULL,
  owner_id bigint REFERENCES users (id) ON DELETE SET NULL,
  filename TEXT NOT NULL,
  storage_path TEXT NOT NULL UNIQUE,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  last_accessed_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  download_count integer NOT NULL DEFAULT 0,
  version integer NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX idx_files_token ON files (token);
