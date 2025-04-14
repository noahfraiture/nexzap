-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE languages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL
);

CREATE TABLE tests (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  content TEXT NOT NULL,
  docker_image VARCHAR(255) NOT NULL
);

CREATE TABLE tutorials (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  language_id UUID NOT NULL REFERENCES languages (id),
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE sheets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  content TEXT NOT NULL,
  test_id UUID NOT NULL REFERENCES tests (id),
  tutorial_id UUID NOT NULL REFERENCES tutorials (id)
);

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  content BYTEA NOT NULL,
  test_id UUID NOT NULL REFERENCES tests (id)
);
