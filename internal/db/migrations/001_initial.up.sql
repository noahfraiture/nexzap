-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tutorials (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title STRING NOT NULL,
  highlight STRING,
  code_editor STRING,
  docker_image STRING NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE sheets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  guide_content TEXT NOT NULL,
  exercise_content TEXT,
  tutorial_id UUID NOT NULL REFERENCES tutorials (id)
);

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  content BYTEA NOT NULL,
  sheet_id UUID NOT NULL REFERENCES sheets (id)
);
