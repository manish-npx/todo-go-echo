DO $$
BEGIN
    -- If legacy "password" exists and "password_hash" doesn't, rename in place.
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'password'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'password_hash'
    ) THEN
        ALTER TABLE users RENAME COLUMN password TO password_hash;
    END IF;

    -- If both columns exist, backfill hash from legacy where needed, then enforce NOT NULL.
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'password'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'password_hash'
    ) THEN
        UPDATE users
        SET password_hash = password
        WHERE password_hash IS NULL AND password IS NOT NULL;
    END IF;

    -- If hash column is missing entirely, create it.
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'password_hash'
    ) THEN
        ALTER TABLE users ADD COLUMN password_hash TEXT;
    END IF;

    -- Final constraint shape expected by application.
    ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
END $$;
