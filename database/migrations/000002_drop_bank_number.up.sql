BEGIN TRANSACTION;

ALTER TABLE "account"
    DROP CONSTRAINT IF EXISTS "account_bank_number_key",
    DROP COLUMN IF EXISTS "bank_number";

COMMIT;
