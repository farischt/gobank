BEGIN TRANSACTION;

DROP EXTENSION IF EXISTS "uuid-ossp";

DROP CONSTRAINT IF EXISTS "session_token_account_id_fkey" ON "session_token";
DROP TABLE IF EXISTS "session_token";

COMMIT;
