BEGIN TRANSACTION;

ALTER TABLE account
    DROP CONSTRAINT account_bank_number_key,
    DROP bank_number;

COMMIT;