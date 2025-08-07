DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.schemata
            WHERE schema_name = 'calendar'
        ) THEN
            CREATE SCHEMA calendar;
            COMMENT ON SCHEMA calendar IS 'Основная схема для приложения';
            REVOKE ALL ON SCHEMA calendar FROM PUBLIC;
            GRANT USAGE ON SCHEMA calendar TO postgres_user;
        END IF;
    END
$$;