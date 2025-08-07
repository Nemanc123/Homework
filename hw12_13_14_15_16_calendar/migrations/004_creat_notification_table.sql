
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'calendar'
              AND table_name = 'notification'
        ) THEN
            CREATE TABLE calendar.notification (
                                    id SERIAL PRIMARY KEY,
                                    title VARCHAR(50) NOT NULL,
                                    date_and_time_of_the_event TIMESTAMP DEFAULT current_timestamp,
                                    id_user BIGINT NOT NULL
            );
            COMMENT ON TABLE calendar.notification IS 'Таблица события';
            COMMENT ON COLUMN calendar.notification.id IS 'id события';
            COMMENT ON COLUMN calendar.notification.title IS 'Заголовок';
            COMMENT ON COLUMN calendar.notification.date_and_time_of_the_event IS 'Дата и время проведения события';
            COMMENT ON COLUMN calendar.notification.id_user IS 'Уникальное id пользователя';
            GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE calendar.notification TO postgres_user;
            CREATE INDEX idx_notification ON calendar.notification(id);
        END IF;
    END
$$