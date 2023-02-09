BEGIN TRANSACTION;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "session_token" (
  "id" text PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "account_id" integer NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

ALTER TABLE "session_token"
    ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
