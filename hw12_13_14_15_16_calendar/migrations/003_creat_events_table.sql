
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.tables
            WHERE table_schema = 'calendar'
              AND table_name = 'events'
        ) THEN
                 CREATE TABLE calendar.events (
                     id SERIAL PRIMARY KEY,
                     title VARCHAR(50) NOT NULL,
                     date_and_time_of_the_event TIMESTAMP DEFAULT current_timestamp,
                     duration_of_the_event TIMESTAMP NOT NULL,
                     description_event VARCHAR(255),
                     id_user BIGINT NOT NULL,
                     time_until_event BIGINT
                 );
            COMMENT ON TABLE calendar.events IS 'Таблица эвентов ';
            COMMENT ON COLUMN calendar.events.id IS 'id события';
            COMMENT ON COLUMN calendar.events.title IS 'Заголовок';
            COMMENT ON COLUMN calendar.events.date_and_time_of_the_event IS 'Дата и время проведения события';
            COMMENT ON COLUMN calendar.events.duration_of_the_event IS 'Продолжительность события';
            COMMENT ON COLUMN calendar.events.description_event IS 'Описание события';
            COMMENT ON COLUMN calendar.events.id_user IS 'Уникальное id пользователя';
            COMMENT ON COLUMN calendar.events.time_until_event IS 'Время отправления уведомления события';
            GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE calendar.events TO postgres_user;
            CREATE INDEX idx_events ON calendar.events(id);
        END IF;
    END
$$