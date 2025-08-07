DO $$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = 'postgres_users') THEN
            CREATE USER postgres_users WITH PASSWORD 'postgres_password';
            RAISE NOTICE 'Role "postgres_users" create successefuly';
        ELSE
            RAISE NOTICE 'Role "postgres_users" already exists';
        END IF;
    END;
$$;

DO $$
    BEGIN
        IF NOT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = 'notify_user') THEN
            CREATE USER notify_user WITH PASSWORD 'postgres_password';
            RAISE NOTICE 'Role "notify_user" create successefuly';
        ELSE
            RAISE NOTICE 'Role "notify_user" already exists';
        END IF;
    END;
$$