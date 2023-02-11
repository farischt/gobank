BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "account" (
  "id" serial PRIMARY KEY,
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR NOT NULL,
  "password" TEXT NOT NULL,
  "bank_number" BIGINT UNIQUE NOT NULL,
  "balance" INTEGER NOT NULL DEFAULT 10,
  "created_at" TIMESTAMP  DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

COMMIT;
