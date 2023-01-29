BEGIN TRANSACTION;

ALTER TABLE transaction 
    DROP CONSTRAINT IF EXISTS transaction_from_id_fkey,
    DROP CONSTRAINT IF EXISTS transaction_to_id_fkey;

DROP TABLE IF EXISTS transaction;

ALTER TABLE account 
    ALTER COLUMN balance TYPE INTEGER,
    ALTER COLUMN balance SET DEFAULT 10,
    ALTER COLUMN balance SET NOT NULL;

COMMIT;