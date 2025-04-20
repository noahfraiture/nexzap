-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tutorials (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  title TEXT NOT NULL,
  highlight TEXT NOT NULL,
  code_editor TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  unlock TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
  UNIQUE (title, version)
);

CREATE TABLE sheets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  tutorial_id UUID NOT NULL REFERENCES tutorials (id),
  page INTEGER NOT NULL,
  guide_content TEXT NOT NULL,
  exercise_content TEXT NOT NULL,
  submission_name TEXT NOT NULL,
  submission_content TEXT NOT NULL,
  docker_image TEXT NOT NULL,
  command TEXT NOT NULL, -- Command to run the test
  UNIQUE (tutorial_id, page)
);

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  name VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  sheet_id UUID NOT NULL REFERENCES sheets (id),
  UNIQUE (name, sheet_id)
);
