CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "pictures" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "url" text NOT NULL,
  "user_id" UUID NOT NULL REFERENCES users("id") ON DELETE CASCADE,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION update_time()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE users SET updated_at = now();
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_time_trigger
AFTER update ON "users"
FOR EACH ROW EXECUTE FUNCTION update_time();